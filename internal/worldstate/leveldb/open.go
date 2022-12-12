// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package leveldb

import (
	"path/filepath"
	"regexp"
	"sync"

	"github.com/hyperledger-labs/orion-server/internal/fileops"
	"github.com/hyperledger-labs/orion-server/internal/worldstate"
	"github.com/hyperledger-labs/orion-server/pkg/logger"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	// underCreationFlag is used to mark that the leveldb
	// instance is being created. If a failure happens during the
	// creation, the retry logic will use this file to
	// detect the partially created store and do cleanup
	// before creating a new levelDB instance
	underCreationFlag = "undercreation"
	// allowedCharsInDBName holds the regexp for allowed characters
	// in a database name
	allowedCharsInDBName = `^[0-9a-zA-Z_\-\.]+$`
)

// LevelDB holds information about all created database
type LevelDB struct {
	dbRootDir   string
	dbs         map[string]*db
	logger      *logger.SugarLogger
	dbsList     sync.RWMutex
	dbNameRegex *regexp.Regexp
	cache       *cache
}

// db - a wrapper on an actual store
type db struct {
	name      string
	file      *leveldb.DB
	mu        sync.RWMutex
	readOpts  *opt.ReadOptions
	writeOpts *opt.WriteOptions
}

var (
	preCreateDBs = append(
		worldstate.SystemDBs(),
		worldstate.DefaultDBName,
	)
)

type Config struct {
	DBRootDir string
	Logger    *logger.SugarLogger
}

// Open opens a leveldb instance to maintain world state
func Open(conf *Config) (*LevelDB, error) {
	exist, err := fileops.Exists(conf.DBRootDir)
	if err != nil {
		return nil, err
	}
	if !exist {
		return openNewLevelDBInstance(conf)
	}

	partialInstanceExist, err := isExistingLevelDBInstanceCreatedPartially(conf.DBRootDir)
	if err != nil {
		return nil, err
	}

	switch {
	case partialInstanceExist:
		if err := fileops.RemoveAll(conf.DBRootDir); err != nil {
			return nil, errors.Wrap(err, "error while removing the existing partially created levelDB instance")
		}

		return openNewLevelDBInstance(conf)
	default:
		return openExistingLevelDBInstance(conf)
	}
}

func isExistingLevelDBInstanceCreatedPartially(dbPath string) (bool, error) {
	empty, err := fileops.IsDirEmpty(dbPath)
	if err != nil {
		return false, err
	}

	if empty {
		return true, nil
	}

	return fileops.Exists(filepath.Join(dbPath, underCreationFlag))
}

func openNewLevelDBInstance(c *Config) (*LevelDB, error) {
	if err := fileops.CreateDir(c.DBRootDir); err != nil {
		return nil, errors.WithMessagef(err, "failed to create director %s", c.DBRootDir)
	}

	underCreationFlagPath := filepath.Join(c.DBRootDir, underCreationFlag)
	if err := fileops.CreateFile(underCreationFlagPath); err != nil {
		return nil, err
	}

	l := &LevelDB{
		dbRootDir:   c.DBRootDir,
		dbs:         make(map[string]*db),
		logger:      c.Logger,
		dbNameRegex: regexp.MustCompile(allowedCharsInDBName),
		cache:       newCache(128),
	}

	for _, dbName := range preCreateDBs {
		if err := l.create(dbName); err != nil {
			return nil, err
		}
	}

	if err := fileops.Remove(underCreationFlagPath); err != nil {
		return nil, errors.WithMessagef(err, "error while removing the under creation flag [%s]", underCreationFlagPath)
	}

	return l, nil
}

func openExistingLevelDBInstance(c *Config) (*LevelDB, error) {
	l := &LevelDB{
		dbRootDir:   c.DBRootDir,
		dbs:         make(map[string]*db),
		logger:      c.Logger,
		dbNameRegex: regexp.MustCompile(allowedCharsInDBName),
		cache:       newCache(128),
	}

	dbNames, err := fileops.ListSubdirs(c.DBRootDir)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to retrieve existing level dbs from %s", c.DBRootDir)
	}

	for _, dbName := range dbNames {
		file, err := leveldb.OpenFile(
			filepath.Join(l.dbRootDir, dbName),
			&opt.Options{ErrorIfMissing: false},
		)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to open leveldb file for database %s", dbName)
		}

		l.dbs[dbName] = &db{
			name:      dbName,
			file:      file,
			readOpts:  &opt.ReadOptions{},
			writeOpts: &opt.WriteOptions{Sync: true},
		}
	}

	return l, nil
}

// Close closes the database instance by closing all leveldb databases
func (l *LevelDB) Close() error {
	l.dbsList.Lock()
	defer l.dbsList.Unlock()

	for name, db := range l.dbs {
		db.mu.Lock()
		defer db.mu.Unlock()

		if err := db.file.Close(); err != nil {
			return errors.Errorf("error while closing database %s, %v", name, err)
		}

		delete(l.dbs, db.name)
	}

	return nil
}

// ValidDBName returns true if the given dbName is valid
func (l *LevelDB) ValidDBName(dbName string) bool {
	return l.dbNameRegex.MatchString(dbName)
}
