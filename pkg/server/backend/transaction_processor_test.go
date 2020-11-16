package backend

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.ibm.com/blockchaindb/protos/types"
	"github.ibm.com/blockchaindb/server/config"
	"github.ibm.com/blockchaindb/server/pkg/blockstore"
	"github.ibm.com/blockchaindb/server/pkg/common/logger"
	"github.ibm.com/blockchaindb/server/pkg/identity"
	"github.ibm.com/blockchaindb/server/pkg/worldstate"
	"github.ibm.com/blockchaindb/server/pkg/worldstate/leveldb"
)

type txProcessorTestEnv struct {
	dbPath         string
	db             *leveldb.LevelDB
	blockStore     *blockstore.Store
	blockStorePath string
	txProcessor    *transactionProcessor
	cleanup        func()
}

func newTxProcessorTestEnv(t *testing.T) *txProcessorTestEnv {
	dir, err := ioutil.TempDir("/tmp", "transactionProcessor")
	require.NoError(t, err)

	c := &logger.Config{
		Level:         "debug",
		OutputPath:    []string{"stdout"},
		ErrOutputPath: []string{"stderr"},
		Encoding:      "console",
	}
	logger, err := logger.New(c)
	require.NoError(t, err)

	dbPath := constructWorldStatePath(dir)
	db, err := leveldb.Open(
		&leveldb.Config{
			DBRootDir: dbPath,
			Logger:    logger,
		},
	)
	if err != nil {
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			t.Errorf("error while removing directory %s, %v", dir, rmErr)
		}
		t.Fatalf("error while creating leveldb, %v", err)
	}

	blockStorePath := constructBlockStorePath(dir)
	blockStore, err := blockstore.Open(
		&blockstore.Config{
			StoreDir: blockStorePath,
			Logger:   logger,
		},
	)
	if err != nil {
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			t.Errorf("error while removing directory %s, %v", dir, rmErr)
		}
		t.Fatalf("error while creating blockstore, %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("error while closing the db instance, %v", err)
		}

		if err := blockStore.Close(); err != nil {
			t.Errorf("error while closing blockstore, %v", err)
		}

		if err := os.RemoveAll(dir); err != nil {
			t.Fatalf("error while removing directory %s, %v", dir, err)
		}
	}

	txProcConf := &txProcessorConfig{
		db:                 db,
		blockStore:         blockStore,
		txQueueLength:      100,
		txBatchQueueLength: 100,
		blockQueueLength:   100,
		maxTxCountPerBatch: 1,
		batchTimeout:       50 * time.Millisecond,
		logger:             logger,
	}
	txProcessor, err := newTransactionProcessor(txProcConf)
	require.NoError(t, err)

	return &txProcessorTestEnv{
		dbPath:         dbPath,
		db:             db,
		blockStorePath: blockStorePath,
		blockStore:     blockStore,
		txProcessor:    txProcessor,
		cleanup:        cleanup,
	}
}

func TestTransactionProcessor(t *testing.T) {
	t.Parallel()

	conf := testConfiguration(t)
	defer os.RemoveAll(conf.Node.Database.LedgerDirectory)

	setup := func(env *txProcessorTestEnv, userID, dbName string) {
		configTx, err := prepareConfigTx(conf)
		require.NoError(t, err)
		require.NoError(t, env.txProcessor.submitTransaction(configTx))

		user := &types.User{
			ID: userID,
			Privilege: &types.Privilege{
				DBPermission: map[string]types.Privilege_Access{
					dbName: types.Privilege_ReadWrite,
				},
			},
		}

		u, err := proto.Marshal(user)
		require.NoError(t, err)

		createUser := []*worldstate.DBUpdates{
			{
				DBName: worldstate.UsersDBName,
				Writes: []*worldstate.KVWithMetadata{
					{
						Key:   string(identity.UserNamespace) + userID,
						Value: u,
						Metadata: &types.Metadata{
							Version: &types.Version{
								BlockNum: 2,
								TxNum:    1,
							},
						},
					},
				},
			},
		}
		require.NoError(t, env.db.Commit(createUser, 2))
		genesisCommitted := func() bool {
			height, _ := env.blockStore.Height()
			return height > uint64(0)
		}
		require.Eventually(t, genesisCommitted, time.Second+5, time.Millisecond*100)
	}

	t.Run("commit a data transaction", func(t *testing.T) {
		t.Parallel()
		env := newTxProcessorTestEnv(t)
		defer env.cleanup()

		setup(env, "testUser", worldstate.DefaultDBName)

		tx := &types.DataTxEnvelope{
			Payload: &types.DataTx{
				UserID:    "testUser",
				DBName:    worldstate.DefaultDBName,
				TxID:      "tx1",
				DataReads: []*types.DataRead{},
				DataWrites: []*types.DataWrite{
					{
						Key:   "test-key1",
						Value: []byte("test-value1"),
					},
				},
			},
		}

		require.NoError(t, env.txProcessor.submitTransaction(tx))

		assertTestKey1InDB := func() bool {
			val, metadata, err := env.db.Get(worldstate.DefaultDBName, "test-key1")
			if err != nil {
				return false
			}
			return bytes.Equal([]byte("test-value1"), val) &&
				proto.Equal(
					&types.Metadata{
						Version: &types.Version{
							BlockNum: 2,
							TxNum:    0,
						},
					},
					metadata,
				)
		}
		require.Eventually(
			t,
			assertTestKey1InDB,
			2*time.Second,
			100*time.Millisecond,
		)

		height, err := env.blockStore.Height()
		require.NoError(t, err)
		require.Equal(t, uint64(2), height)

		genesisHash, err := env.blockStore.GetHash(1)
		require.NoError(t, err)
		require.NotNil(t, genesisHash)
		genesisHashBase, err := env.blockStore.GetBaseHeaderHash(1)
		require.NoError(t, err)
		require.NotNil(t, genesisHashBase)

		expectedBlock := &types.Block{
			Header: &types.BlockHeader{
				BaseHeader: &types.BlockHeaderBase{
					Number:                 2,
					PreviousBaseHeaderHash: genesisHashBase,
					LastCommittedBlockHash: genesisHash,
					LastCommittedBlockNum:  1,
				},
				SkipchainHashes: [][]byte{genesisHash},
				ValidationInfo: []*types.ValidationInfo{
					{
						Flag: types.Flag_VALID,
					},
				},
			},
			Payload: &types.Block_DataTxEnvelopes{
				DataTxEnvelopes: &types.DataTxEnvelopes{
					Envelopes: []*types.DataTxEnvelope{
						tx,
					},
				},
			},
		}

		block, err := env.blockStore.Get(2)
		require.NoError(t, err)
		require.True(t, proto.Equal(expectedBlock, block))
	})
}

func testConfiguration(t *testing.T) *config.Configurations {
	ledgerDir, err := ioutil.TempDir("/tmp", "server")
	require.NoError(t, err)

	return &config.Configurations{
		Node: config.NodeConf{
			Identity: config.IdentityConf{
				ID:              "bdb-node-1",
				CertificatePath: "./testdata/node.cert",
				KeyPath:         "./testdata/node.key",
			},
			Network: config.NetworkConf{
				Address: "127.0.0.1",
				Port:    0,
			},
			Database: config.DatabaseConf{
				Name:            "leveldb",
				LedgerDirectory: ledgerDir,
			},
			QueueLength: config.QueueLengthConf{
				Transaction:               1000,
				ReorderedTransactionBatch: 100,
				Block:                     100,
			},
			LogLevel: "debug",
		},
		Consensus: config.ConsensusConf{
			Algorithm:                   "raft",
			MaxBlockSize:                2,
			MaxTransactionCountPerBlock: 1,
			BlockTimeout:                50 * time.Millisecond,
		},
		Admin: config.AdminConf{
			ID:              "admin",
			CertificatePath: "./testdata/admin.cert",
		},
		RootCA: config.RootCAConf{
			CertificatePath: "./testdata/rootca.cert",
		},
	}
}
