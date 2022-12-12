// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package leveldb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger-labs/orion-server/internal/worldstate"
	"github.com/hyperledger-labs/orion-server/pkg/types"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	lastCommittedBlockNumberKey = []byte("lastCommittedBlockNumber")
)

// Exist returns true if the given database exist. Otherwise, it returns false.
func (l *LevelDB) Exist(dbName string) bool {
	_, ok := l.GetDB(dbName)
	return ok
}

// ListDBs list all user databases
func (l *LevelDB) ListDBs() []string {
	dbsToExclude := make(map[string]struct{})
	for _, name := range preCreateDBs {
		dbsToExclude[name] = struct{}{}
	}

	var dbNames []string
	l.Range(func(key, value interface{}) bool {
		name := key.(string)
		if _, ok := dbsToExclude[name]; !ok {
			dbNames = append(dbNames, name)
		}
		return true
	})

	return dbNames
}

// Height returns the block height of the state database. In other words, it
// returns the last committed block number
func (l *LevelDB) Height() (uint64, error) {
	db, ok := l.GetDB(worldstate.MetadataDBName)
	if !ok {
		return 0, errors.Errorf("unable to retrieve the state database height due to missing metadataDB")
	}

	blockNumberEnc, err := db.file.Get(lastCommittedBlockNumberKey, &opt.ReadOptions{})
	if err != nil && err != leveldb.ErrNotFound {
		return 0, errors.Wrap(err, "error while retrieving the state database height")
	}

	if err == leveldb.ErrNotFound {
		return 0, nil
	}

	blockNumberDec, err := binary.ReadUvarint(bytes.NewBuffer(blockNumberEnc))
	if err != nil {
		return 0, errors.Wrap(err, "error while decoding the stored height")
	}

	return blockNumberDec, nil
}

// Get returns the value of the key present in the database.
func (l *LevelDB) Get(dbName string, key string) ([]byte, *types.Metadata, error) {
	db, ok := l.GetDB(dbName)
	if !ok {
		return nil, nil, &DBNotFoundErr{
			dbName: dbName,
		}
	}

	peristed, err := l.cache.getState(dbName, key)
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "failed to retrieve leveldb key [%s] from database %s through cache", key, dbName)
	}
	if peristed != nil {
		return peristed.Value, peristed.Metadata, nil
	}

	var dbval []byte
	if db.snap != nil {
		dbval, err = db.snap.Get([]byte(key), db.readOpts)
	} else {
		dbval, err = db.file.Get([]byte(key), db.readOpts)
	}
	if err == leveldb.ErrNotFound {
		if err = l.cache.putState(dbName, key, nil); err != nil {
			return nil, nil, err
		}
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "failed to retrieve leveldb key [%s] from database %s", key, dbName)
	}

	if err = l.cache.putState(dbName, key, dbval); err != nil {
		return nil, nil, err
	}

	persisted := &types.ValueWithMetadata{}
	if err := proto.Unmarshal(dbval, persisted); err != nil {
		return nil, nil, err
	}

	return persisted.Value, persisted.Metadata, nil
}

// GetVersion returns the version of the key present in the database
func (l *LevelDB) GetVersion(dbName string, key string) (*types.Version, error) {
	_, metadata, err := l.Get(dbName, key)
	if err != nil {
		return nil, err
	}

	return metadata.GetVersion(), nil
}

// GetACL returns the access control rule for the given key present in the database
func (l *LevelDB) GetACL(dbName, key string) (*types.AccessControl, error) {
	_, metadata, err := l.Get(dbName, key)
	if err != nil {
		return nil, err
	}

	return metadata.GetAccessControl(), nil
}

// Has returns true if the key exist in the database
func (l *LevelDB) Has(dbName, key string) (bool, error) {
	db, ok := l.GetDB(dbName)
	if !ok {
		return false, nil
	}

	return db.file.Has([]byte(key), nil)
}

// GetConfig returns the cluster configuration
func (l *LevelDB) GetConfig() (*types.ClusterConfig, *types.Metadata, error) {
	configSerialized, metadata, err := l.Get(worldstate.ConfigDBName, worldstate.ConfigKey)
	if err != nil {
		return nil, nil, err
	}

	config := &types.ClusterConfig{}
	if err := proto.Unmarshal(configSerialized, config); err != nil {
		return nil, nil, errors.Wrap(err, "error while unmarshaling committed cluster configuration")
	}

	return config, metadata, nil
}

// GetIndexDefinition returns the index definition of a given database
func (l *LevelDB) GetIndexDefinition(dbName string) ([]byte, *types.Metadata, error) {
	return l.Get(worldstate.DatabasesDBName, dbName)
}

// GetIterator returns an iterator to fetch values associated with a range of keys
// startKey is inclusive while the endKey is exclusive. An empty startKey (i.e., "") denotes that
// the caller wants from the first key in the database (lexicographic order). An empty
// endKey (i.e., "") denotes that the caller wants till the last key in the database (lexicographic order).
func (l *LevelDB) GetIterator(dbName string, startKey, endKey string) (worldstate.Iterator, error) {
	db, ok := l.GetDB(dbName)

	if !ok || db == nil {
		l.logger.Errorf("database %s does not exist", dbName)
		return nil, errors.Errorf("database %s does not exist", dbName)
	}

	r := &util.Range{}
	if startKey == "" {
		r.Start = nil
	} else {
		r.Start = []byte(startKey)
	}

	if endKey == "" {
		r.Limit = nil
	} else {
		r.Limit = []byte(endKey)
	}

	return db.file.NewIterator(r, &opt.ReadOptions{}), nil
}

// Commit commits the updates to the database
func (l *LevelDB) Commit(dbsUpdates map[string]*worldstate.DBUpdates, blockNumber uint64) error {
	for dbName, updates := range dbsUpdates {
		db, ok := l.GetDB(dbName)

		if !ok || db == nil {
			l.logger.Errorf("database %s does not exist", dbName)
			return errors.Errorf("database %s does not exist", dbName)
		}

		start := time.Now()
		if err := l.commitToDB(dbName, db, updates); err != nil {
			return err
		}
		snap, err := db.file.GetSnapshot()
		if err != nil {
			return err
		}
		prev := db.snap
		db.snap = snap
		if prev != nil {
			prev.Release()
		}
		l.logger.Debugf("changes committed to the database %s, took %d ms, available dbs are [%s]", dbName, time.Since(start).Milliseconds(), "")
	}

	db, exists := l.GetDB(worldstate.MetadataDBName)
	if !exists {
		l.logger.Errorf("metadata database does not exist, available dbs are [%+v]", "")
		return errors.Errorf("metadata database does not exist")
	}

	b := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(b, blockNumber)
	if err := db.file.Put(lastCommittedBlockNumberKey, b, &opt.WriteOptions{}); err != nil {
		return errors.Wrapf(err, "error while storing the last committed block number [%d] to the metadataDB", blockNumber)
	}

	return nil
}

func (l *LevelDB) commitToDB(dbName string, db *Db, updates *worldstate.DBUpdates) error {
	batch := &leveldb.Batch{}

	for _, kv := range updates.Writes {
		dbval, err := proto.Marshal(
			&types.ValueWithMetadata{
				Value:    kv.Value,
				Metadata: kv.Metadata,
			},
		)
		if err != nil {
			return errors.WithMessagef(err, "failed to marshal the constructed dbValue [%v]", kv.Value)
		}

		batch.Put([]byte(kv.Key), dbval)
		l.cache.putStateIfExist(dbName, kv.Key, dbval)
	}

	for _, key := range updates.Deletes {
		batch.Delete([]byte(key))
		l.cache.delState(dbName, key)
	}

	if err := db.file.Write(batch, db.writeOpts); err != nil {
		return errors.Wrapf(err, "error while writing an update batch to database [%s]", db.name)
	}

	if dbName != worldstate.DatabasesDBName {
		return nil
	}

	// if node fails during the creation or deletion of
	// databases, during the recovery, these operations
	// will be repeated again. Given that create() and
	// delete() are a no-op when the db exist and not-exist,
	// respectively, we don't need anything special to
	// handle failures

	// we also assume the union of dbNames in create
	// and delete list to be unique which is to be ensured
	// by the validator.

	for _, kv := range updates.Writes {
		dbName := kv.Key
		if err := l.create(dbName); err != nil {
			return err
		}
	}

	for _, dbName := range updates.Deletes {
		if err := l.delete(dbName); err != nil {
			return err
		}
	}

	return nil
}

// create creates a database. It does not return an error when the database already exist.
func (l *LevelDB) create(dbName string) error {
	file, err := leveldb.OpenFile(filepath.Join(l.dbRootDir, dbName), &opt.Options{})
	if err != nil {
		return errors.WithMessagef(err, "failed to open leveldb file for database %s", dbName)
	}

	success := l.SetDB(dbName, &Db{
		name:      dbName,
		file:      file,
		readOpts:  &opt.ReadOptions{},
		writeOpts: &opt.WriteOptions{Sync: true},
	})

	if !success {
		_ = file.Close()
		l.logger.Debugf("Skipping %s cause database already exists", dbName)
		return nil
	}

	return nil
}

// delete deletes a database. It does not return an error when the database does not exist.
// delete would be called only by the Commit() when processing delete entries associated with
// the _db
func (l *LevelDB) delete(dbName string) error {
	db, existed := l.GetAndDelDB(dbName)
	if !existed {
		return nil
	}

	if err := db.file.Close(); err != nil {
		return errors.Wrapf(err, "error while closing the database [%s] before delete", dbName)
	}

	if err := os.RemoveAll(filepath.Join(l.dbRootDir, dbName)); err != nil {
		return errors.Wrapf(err, "error while deleting database [%s]", dbName)
	}

	return nil
}

// DBNotFoundErr denotes that the given dbName is not present in the database
type DBNotFoundErr struct {
	dbName string
}

func (e *DBNotFoundErr) Error() string {
	return fmt.Sprintf("database %s does not exist", e.dbName)
}
