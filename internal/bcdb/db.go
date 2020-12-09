package backend

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.ibm.com/blockchaindb/server/config"
	"github.ibm.com/blockchaindb/server/internal/blockstore"
	"github.ibm.com/blockchaindb/server/internal/fileops"
	"github.ibm.com/blockchaindb/server/internal/identity"
	"github.ibm.com/blockchaindb/server/internal/provenance"
	"github.ibm.com/blockchaindb/server/internal/worldstate"
	"github.ibm.com/blockchaindb/server/internal/worldstate/leveldb"
	"github.ibm.com/blockchaindb/server/pkg/crypto"
	"github.ibm.com/blockchaindb/server/pkg/logger"
	"github.ibm.com/blockchaindb/server/pkg/types"
)

//go:generate mockery --dir . --name DB --case underscore --output mocks/

// DB encapsulates functionality required to operate with database state
type DB interface {
	// LedgerHeight returns current height of the ledger
	LedgerHeight() (uint64, error)

	// Height returns ledger height
	Height() (uint64, error)

	// DoesUserExist checks whenever user with given userID exists
	DoesUserExist(userID string) (bool, error)

	// GetCertificate returns the certificate associated with useID, if it exists.
	GetCertificate(userID string) (*x509.Certificate, error)

	// GetUser retrieves user' record
	GetUser(querierUserID, targetUserID string) (*types.GetUserResponseEnvelope, error)

	// GetConfig returns database configuration
	GetConfig() (*types.GetConfigResponseEnvelope, error)

	// GetNodeConfig returns single node subsection of database configuration
	GetNodeConfig(nodeID string) (*types.GetNodeConfigResponseEnvelope, error)

	// GetDBStatus returns status for database, checks whenever database was created
	GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error)

	// GetData retrieves values for given key
	GetData(dbName, querierUserID, key string) (*types.GetDataResponseEnvelope, error)

	// GetBlockHeader returns ledger block header
	GetBlockHeader(userID string, blockNum uint64) (*types.GetBlockResponseEnvelope, error)

	// GetTxProof returns intermediate hashes to recalculate merkle tree root from tx hash
	GetTxProof(userID string, blockNum uint64, txIdx uint64) (*types.GetTxProofResponseEnvelope, error)

	// GetLedgerPath returns list of blocks that forms shortest path in skip list chain in ledger
	GetLedgerPath(userID string, start, end uint64) (*types.GetLedgerPathResponseEnvelope, error)

	// GetValues returns all values associated with a given key
	GetValues(dbName, key string) (*types.GetHistoricalDataResponseEnvelope, error)

	// GetValueAt returns the value of a given key at a particular version
	GetValueAt(dbName, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error)

	// GetPreviousValues returns previous values of a given key and a version. The number of records returned would be limited
	// by the limit parameters.
	GetPreviousValues(dbname, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error)

	// GetNextValues returns next values of a given key and a version. The number of records returned would be limited
	// by the limit parameters.
	GetNextValues(dbname, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error)

	// GetValuesReadByUser returns all values read by a given user
	GetValuesReadByUser(userID string) (*types.GetDataReadByResponseEnvelope, error)

	// GetValuesReadByUser returns all values read by a given user
	GetValuesWrittenByUser(userID string) (*types.GetDataWrittenByResponseEnvelope, error)

	// GetReaders returns all userIDs who have accessed a given key as well as the access frequency
	GetReaders(dbName, key string) (*types.GetDataReadersResponseEnvelope, error)

	// GetReaders returns all userIDs who have accessed a given key as well as the access frequency
	GetWriters(dbName, key string) (*types.GetDataWritersResponseEnvelope, error)

	// GetTxIDsSubmittedByUser returns all ids of all transactions submitted by a given user
	GetTxIDsSubmittedByUser(userID string) (*types.GetTxIDsSubmittedByResponseEnvelope, error)

	// GetTxReceipt returns transaction receipt - block header of ledger block that contains the transaction
	// and transaction index inside the block
	GetTxReceipt(userId string, txID string) (*types.GetTxReceiptResponseEnvelope, error)

	// SubmitTransaction submits transaction to the database
	SubmitTransaction(tx interface{}) error

	// BootstrapDB given bootstrap configuration initialize database by
	// creating required system tables to include database meta data
	BootstrapDB(conf *config.Configurations) error

	// IsReady returns true once instance of the DB is properly initiated, meaning
	// all system tables was created successfully
	IsReady() (bool, error)

	// IsDBExists returns true if database with given name is exists otherwise false
	IsDBExists(name string) bool

	// Close frees and closes resources allocated by database instance
	Close() error
}

type db struct {
	worldstateQueryProcessor *worldstateQueryProcessor
	ledgerQueryProcessor     *ledgerQueryProcessor
	provenanceQueryProcessor *provenanceQueryProcessor
	txProcessor              *transactionProcessor
	db                       worldstate.DB
	blockStore               *blockstore.Store
	provenanceStore          *provenance.Store
	logger                   *logger.SugarLogger
}

// NewDB creates a new database backend which handles both the queries and transactions.
func NewDB(conf *config.Configurations, logger *logger.SugarLogger) (DB, error) {
	if conf.Node.Database.Name != "leveldb" {
		return nil, errors.New("only leveldb is supported as the state database")
	}

	ledgerDir := conf.Node.Database.LedgerDirectory
	if err := createLedgerDir(ledgerDir); err != nil {
		return nil, err
	}

	levelDB, err := leveldb.Open(
		&leveldb.Config{
			DBRootDir: constructWorldStatePath(ledgerDir),
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating the world state database")
	}

	blockStore, err := blockstore.Open(
		&blockstore.Config{
			StoreDir: constructBlockStorePath(ledgerDir),
			Logger:   logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating the block store")
	}

	provenanceStore, err := provenance.Open(
		&provenance.Config{
			StoreDir: constructProvenanceStorePath(ledgerDir),
			Logger:   logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating the block sotre")
	}

	querier := identity.NewQuerier(levelDB)

	signer, err := crypto.NewSigner(&crypto.SignerOptions{KeyFilePath: conf.Node.Identity.KeyPath})
	if err != nil {
		return nil, errors.Wrap(err, "can't load private key")
	}

	worldstateQueryProcessor := newWorldstateQueryProcessor(
		&worldstateQueryProcessorConfig{
			nodeID:          conf.Node.Identity.ID,
			signer:          signer,
			db:              levelDB,
			blockStore:      blockStore,
			identityQuerier: querier,
			logger:          logger,
		},
	)

	ledgerQueryProcessorConfig := &ledgerQueryProcessorConfig{
		nodeID:          conf.Node.Identity.ID,
		signer:          signer,
		db:              levelDB,
		blockStore:      blockStore,
		provenanceStore: provenanceStore,
		identityQuerier: querier,
		logger:          logger,
	}
	ledgerQueryProcessor := newLedgerQueryProcessor(ledgerQueryProcessorConfig)


	provenanceQueryProcessor := newProvenanceQueryProcessor(
		&provenanceQueryProcessorConfig{
			nodeID:          conf.Node.Identity.ID,
			signer:          signer,
			provenanceStore: provenanceStore,
			logger:          logger,
		},
	)

	txProcessor, err := newTransactionProcessor(
		&txProcessorConfig{
			db:                 levelDB,
			blockStore:         blockStore,
			provenanceStore:    provenanceStore,
			txQueueLength:      conf.Node.QueueLength.Transaction,
			txBatchQueueLength: conf.Node.QueueLength.ReorderedTransactionBatch,
			blockQueueLength:   conf.Node.QueueLength.Block,
			maxTxCountPerBatch: conf.Consensus.MaxTransactionCountPerBlock,
			batchTimeout:       conf.Consensus.BlockTimeout,
			logger:             logger,
		},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "can't initiate tx processor")
	}

	return &db{
		worldstateQueryProcessor: worldstateQueryProcessor,
		ledgerQueryProcessor:     ledgerQueryProcessor,
		provenanceQueryProcessor: provenanceQueryProcessor,
		txProcessor:              txProcessor,
		db:                       levelDB,
		blockStore:               blockStore,
		provenanceStore:          provenanceStore,
		logger:                   logger,
	}, nil
}

// BootstrapDB bootstraps DB with system tables
func (d *db) BootstrapDB(conf *config.Configurations) error {
	configTx, err := prepareConfigTx(conf)
	if err != nil {
		return errors.Wrap(err, "failed to prepare and commit a configuration transaction")
	}

	if err := d.txProcessor.submitTransaction(configTx); err != nil {
		return errors.Wrap(err, "error while committing configuration transaction")
	}
	return nil
}

// IsReady returns true once instance of the DB is properly initiated, meaning
// all system tables was created successfully
func (d *db) IsReady() (bool, error) {
	height, err := d.LedgerHeight()
	if err != nil {
		return false, err
	}

	dbHeight, err := d.Height()
	if err != nil {
		return false, err
	}

	for _, sysDB := range worldstate.SystemDBs() {
		if !d.IsDBExists(sysDB) {
			return false, nil
		}
	}
	return height > 0 && dbHeight > 0, nil
}

// LedgerHeight returns ledger height
func (d *db) LedgerHeight() (uint64, error) {
	return d.worldstateQueryProcessor.blockStore.Height()
}

// Height returns ledger height
func (d *db) Height() (uint64, error) {
	return d.worldstateQueryProcessor.db.Height()
}

// DoesUserExist checks whenever userID exists
func (d *db) DoesUserExist(userID string) (bool, error) {
	return d.worldstateQueryProcessor.identityQuerier.DoesUserExist(userID)
}

func (d *db) GetCertificate(userID string) (*x509.Certificate, error) {
	return d.worldstateQueryProcessor.identityQuerier.GetCertificate(userID)
}

// GetUser returns user's record
func (d *db) GetUser(querierUserID, targetUserID string) (*types.GetUserResponseEnvelope, error) {
	return d.worldstateQueryProcessor.getUser(querierUserID, targetUserID)
}

// GetNodeConfig returns single node subsection of database configuration
func (d *db) GetNodeConfig(nodeID string) (*types.GetNodeConfigResponseEnvelope, error) {
	return d.worldstateQueryProcessor.getNodeConfig(nodeID)
}

// GetConfig returns database configuration
func (d *db) GetConfig() (*types.GetConfigResponseEnvelope, error) {
	return d.worldstateQueryProcessor.getConfig()
}

// GetDBStatus returns database status
func (d *db) GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error) {
	return d.worldstateQueryProcessor.getDBStatus(dbName)
}

// SubmitTransaction submits transaction
func (d *db) SubmitTransaction(tx interface{}) error {
	return d.txProcessor.submitTransaction(tx)
}

// GetData returns value for provided key
func (d *db) GetData(dbName, querierUserID, key string) (*types.GetDataResponseEnvelope, error) {
	return d.worldstateQueryProcessor.getData(dbName, querierUserID, key)
}

func (d *db) IsDBExists(name string) bool {
	return d.worldstateQueryProcessor.isDBExists(name)
}

func (d *db) GetBlockHeader(userID string, blockNum uint64) (*types.GetBlockResponseEnvelope, error) {
	return d.ledgerQueryProcessor.getBlockHeader(userID, blockNum)
}

func (d *db) GetTxProof(userID string, blockNum uint64, txIdx uint64) (*types.GetTxProofResponseEnvelope, error) {
	return d.ledgerQueryProcessor.getProof(userID, blockNum, txIdx)
}

func (d *db) GetLedgerPath(userID string, start, end uint64) (*types.GetLedgerPathResponseEnvelope, error) {
	return d.ledgerQueryProcessor.getPath(userID, start, end)
}

func (d *db) GetTxReceipt(userId string, txID string) (*types.GetTxReceiptResponseEnvelope, error) {
	return d.ledgerQueryProcessor.getTxReceipt(userId, txID)
}

// GetValues returns all values associated with a given key
func (d *db) GetValues(dbName, key string) (*types.GetHistoricalDataResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetValues(dbName, key)
}

// GetValueAt returns the value of a given key at a particular version
func (d *db) GetValueAt(dbName, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetValueAt(dbName, key, version)
}

// GetPreviousValues returns previous values of a given key and a version. The number of records returned would be limited
// by the limit parameters.
func (d *db) GetPreviousValues(dbName, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetPreviousValues(dbName, key, version)
}

// GetNextValues returns next values of a given key and a version. The number of records returned would be limited
// by the limit parameters.
func (d *db) GetNextValues(dbName, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetNextValues(dbName, key, version)
}

// GetValuesReadByUser returns all values read by a given user
func (d *db) GetValuesReadByUser(userID string) (*types.GetDataReadByResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetValuesReadByUser(userID)
}

// GetValuesReadByUser returns all values read by a given user
func (d *db) GetValuesWrittenByUser(userID string) (*types.GetDataWrittenByResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetValuesWrittenByUser(userID)
}

// GetReaders returns all userIDs who have accessed a given key as well as the access frequency
func (d *db) GetReaders(dbName, key string) (*types.GetDataReadersResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetReaders(dbName, key)
}

// GetReaders returns all userIDs who have accessed a given key as well as the access frequency
func (d *db) GetWriters(dbName, key string) (*types.GetDataWritersResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetWriters(dbName, key)
}

// GetTxIDsSubmittedByUser returns all ids of all transactions submitted by a given user
func (d *db) GetTxIDsSubmittedByUser(userID string) (*types.GetTxIDsSubmittedByResponseEnvelope, error) {
	return d.provenanceQueryProcessor.GetTxIDsSubmittedByUser(userID)
}

// Close closes and release resources used by db
func (d *db) Close() error {
	if err := d.txProcessor.close(); err != nil {
		return errors.WithMessage(err, "error while closing the transaction processor")
	}

	if err := d.db.Close(); err != nil {
		return errors.WithMessage(err, "error while closing the worldstate database")
	}

	if err := d.provenanceStore.Close(); err != nil {
		return errors.WithMessage(err, "error while closing the provenance store")
	}

	if err := d.blockStore.Close(); err != nil {
		return errors.WithMessage(err, "error while closing the block store")
	}

	return nil
}

func prepareConfigTx(conf *config.Configurations) (*types.ConfigTxEnvelope, error) {
	certs, err := readCerts(conf)
	if err != nil {
		return nil, err
	}

	clusterConfig := &types.ClusterConfig{
		Nodes: []*types.NodeConfig{
			{
				ID:          conf.Node.Identity.ID,
				Certificate: certs.nodeCert,
				Address:     conf.Node.Network.Address,
				Port:        conf.Node.Network.Port,
			},
		},
		Admins: []*types.Admin{
			{
				ID:          conf.Admin.ID,
				Certificate: certs.adminCert,
			},
		},
		RootCACertificate: certs.rootCACert,
	}

	return &types.ConfigTxEnvelope{
		Payload: &types.ConfigTx{
			TxID:      uuid.New().String(), // TODO: we need to change TxID to string
			NewConfig: clusterConfig,
		},
		// TODO: we can make the node itself sign the transaction
	}, nil
}

type certsInGenesisConfig struct {
	nodeCert   []byte
	adminCert  []byte
	rootCACert []byte
}

func readCerts(conf *config.Configurations) (*certsInGenesisConfig, error) {
	nodeCert, err := ioutil.ReadFile(conf.Node.Identity.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading node certificate %s", conf.Node.Identity.CertificatePath)
	}
	nodePemCert, _ := pem.Decode(nodeCert)

	adminCert, err := ioutil.ReadFile(conf.Admin.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading admin certificate %s", conf.Admin.CertificatePath)
	}
	adminPemCert, _ := pem.Decode(adminCert)

	rootCACert, err := ioutil.ReadFile(conf.RootCA.CertificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error while reading rootCA certificate %s", conf.RootCA.CertificatePath)
	}
	rootCAPemCert, _ := pem.Decode(rootCACert)

	return &certsInGenesisConfig{
		nodeCert:   nodePemCert.Bytes,
		adminCert:  adminPemCert.Bytes,
		rootCACert: rootCAPemCert.Bytes,
	}, nil
}

func createLedgerDir(dir string) error {
	exist, err := fileops.Exists(dir)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	return fileops.CreateDir(dir)
}