// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	context "context"

	errors "github.com/hyperledger-labs/orion-server/internal/errors"
	mock "github.com/stretchr/testify/mock"

	time "time"

	types "github.com/hyperledger-labs/orion-server/pkg/types"

	x509 "crypto/x509"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *DB) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataQuery provides a mock function with given fields: ctx, dbName, querierUserID, query
func (_m *DB) DataQuery(ctx context.Context, dbName string, querierUserID string, query []byte) (*types.DataQueryResponseEnvelope, error) {
	ret := _m.Called(ctx, dbName, querierUserID, query)

	var r0 *types.DataQueryResponseEnvelope
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []byte) *types.DataQueryResponseEnvelope); ok {
		r0 = rf(ctx, dbName, querierUserID, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.DataQueryResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, []byte) error); ok {
		r1 = rf(ctx, dbName, querierUserID, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DoesUserExist provides a mock function with given fields: userID
func (_m *DB) DoesUserExist(userID string) (bool, error) {
	ret := _m.Called(userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAugmentedBlockHeader provides a mock function with given fields: userID, blockNum
func (_m *DB) GetAugmentedBlockHeader(userID string, blockNum uint64) (*types.GetAugmentedBlockHeaderResponseEnvelope, error) {
	ret := _m.Called(userID, blockNum)

	var r0 *types.GetAugmentedBlockHeaderResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64) *types.GetAugmentedBlockHeaderResponseEnvelope); ok {
		r0 = rf(userID, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetAugmentedBlockHeaderResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(userID, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeader provides a mock function with given fields: userID, blockNum
func (_m *DB) GetBlockHeader(userID string, blockNum uint64) (*types.GetBlockResponseEnvelope, error) {
	ret := _m.Called(userID, blockNum)

	var r0 *types.GetBlockResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64) *types.GetBlockResponseEnvelope); ok {
		r0 = rf(userID, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetBlockResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(userID, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCertificate provides a mock function with given fields: userID
func (_m *DB) GetCertificate(userID string) (*x509.Certificate, error) {
	ret := _m.Called(userID)

	var r0 *x509.Certificate
	if rf, ok := ret.Get(0).(func(string) *x509.Certificate); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*x509.Certificate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClusterStatus provides a mock function with given fields: noCerts
func (_m *DB) GetClusterStatus(noCerts bool) (*types.GetClusterStatusResponseEnvelope, error) {
	ret := _m.Called(noCerts)

	var r0 *types.GetClusterStatusResponseEnvelope
	if rf, ok := ret.Get(0).(func(bool) *types.GetClusterStatusResponseEnvelope); ok {
		r0 = rf(noCerts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetClusterStatusResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(noCerts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields: querierUserID
func (_m *DB) GetConfig(querierUserID string) (*types.GetConfigResponseEnvelope, error) {
	ret := _m.Called(querierUserID)

	var r0 *types.GetConfigResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetConfigResponseEnvelope); ok {
		r0 = rf(querierUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetConfigResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(querierUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfigBlock provides a mock function with given fields: querierUserID, blockNumber
func (_m *DB) GetConfigBlock(querierUserID string, blockNumber uint64) (*types.GetConfigBlockResponseEnvelope, error) {
	ret := _m.Called(querierUserID, blockNumber)

	var r0 *types.GetConfigBlockResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64) *types.GetConfigBlockResponseEnvelope); ok {
		r0 = rf(querierUserID, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetConfigBlockResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(querierUserID, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDBStatus provides a mock function with given fields: dbName
func (_m *DB) GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error) {
	ret := _m.Called(dbName)

	var r0 *types.GetDBStatusResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetDBStatusResponseEnvelope); ok {
		r0 = rf(dbName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDBStatusResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dbName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetData provides a mock function with given fields: dbName, querierUserID, key
func (_m *DB) GetData(dbName string, querierUserID string, key string) (*types.GetDataResponseEnvelope, error) {
	ret := _m.Called(dbName, querierUserID, key)

	var r0 *types.GetDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, string) *types.GetDataResponseEnvelope); ok {
		r0 = rf(dbName, querierUserID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(dbName, querierUserID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDataProof provides a mock function with given fields: userID, blockNum, dbname, key, deleted
func (_m *DB) GetDataProof(userID string, blockNum uint64, dbname string, key string, deleted bool) (*types.GetDataProofResponseEnvelope, error) {
	ret := _m.Called(userID, blockNum, dbname, key, deleted)

	var r0 *types.GetDataProofResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64, string, string, bool) *types.GetDataProofResponseEnvelope); ok {
		r0 = rf(userID, blockNum, dbname, key, deleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataProofResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, string, string, bool) error); ok {
		r1 = rf(userID, blockNum, dbname, key, deleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeletedValues provides a mock function with given fields: dbname, key
func (_m *DB) GetDeletedValues(dbname string, key string) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbname, key)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbname, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dbname, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLedgerPath provides a mock function with given fields: userID, start, end
func (_m *DB) GetLedgerPath(userID string, start uint64, end uint64) (*types.GetLedgerPathResponseEnvelope, error) {
	ret := _m.Called(userID, start, end)

	var r0 *types.GetLedgerPathResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64, uint64) *types.GetLedgerPathResponseEnvelope); ok {
		r0 = rf(userID, start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetLedgerPathResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, uint64) error); ok {
		r1 = rf(userID, start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMostRecentValueAtOrBelow provides a mock function with given fields: dbName, key, version
func (_m *DB) GetMostRecentValueAtOrBelow(dbName string, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbName, key, version)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, *types.Version) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbName, key, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *types.Version) error); ok {
		r1 = rf(dbName, key, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNextValues provides a mock function with given fields: dbname, key, version
func (_m *DB) GetNextValues(dbname string, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbname, key, version)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, *types.Version) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbname, key, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *types.Version) error); ok {
		r1 = rf(dbname, key, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNodeConfig provides a mock function with given fields: nodeID
func (_m *DB) GetNodeConfig(nodeID string) (*types.GetNodeConfigResponseEnvelope, error) {
	ret := _m.Called(nodeID)

	var r0 *types.GetNodeConfigResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetNodeConfigResponseEnvelope); ok {
		r0 = rf(nodeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetNodeConfigResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(nodeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPreviousValues provides a mock function with given fields: dbname, key, version
func (_m *DB) GetPreviousValues(dbname string, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbname, key, version)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, *types.Version) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbname, key, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *types.Version) error); ok {
		r1 = rf(dbname, key, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReaders provides a mock function with given fields: dbName, key
func (_m *DB) GetReaders(dbName string, key string) (*types.GetDataReadersResponseEnvelope, error) {
	ret := _m.Called(dbName, key)

	var r0 *types.GetDataReadersResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetDataReadersResponseEnvelope); ok {
		r0 = rf(dbName, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataReadersResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dbName, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxIDsSubmittedByUser provides a mock function with given fields: userID
func (_m *DB) GetTxIDsSubmittedByUser(userID string) (*types.GetTxIDsSubmittedByResponseEnvelope, error) {
	ret := _m.Called(userID)

	var r0 *types.GetTxIDsSubmittedByResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetTxIDsSubmittedByResponseEnvelope); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetTxIDsSubmittedByResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxProof provides a mock function with given fields: userID, blockNum, txIdx
func (_m *DB) GetTxProof(userID string, blockNum uint64, txIdx uint64) (*types.GetTxProofResponseEnvelope, error) {
	ret := _m.Called(userID, blockNum, txIdx)

	var r0 *types.GetTxProofResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64, uint64) *types.GetTxProofResponseEnvelope); ok {
		r0 = rf(userID, blockNum, txIdx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetTxProofResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, uint64) error); ok {
		r1 = rf(userID, blockNum, txIdx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxReceipt provides a mock function with given fields: userId, txID
func (_m *DB) GetTxReceipt(userId string, txID string) (*types.TxReceiptResponseEnvelope, error) {
	ret := _m.Called(userId, txID)

	var r0 *types.TxReceiptResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.TxReceiptResponseEnvelope); ok {
		r0 = rf(userId, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.TxReceiptResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userId, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: querierUserID, targetUserID
func (_m *DB) GetUser(querierUserID string, targetUserID string) (*types.GetUserResponseEnvelope, error) {
	ret := _m.Called(querierUserID, targetUserID)

	var r0 *types.GetUserResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetUserResponseEnvelope); ok {
		r0 = rf(querierUserID, targetUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetUserResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(querierUserID, targetUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValueAt provides a mock function with given fields: dbName, key, version
func (_m *DB) GetValueAt(dbName string, key string, version *types.Version) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbName, key, version)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, *types.Version) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbName, key, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *types.Version) error); ok {
		r1 = rf(dbName, key, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValues provides a mock function with given fields: dbName, key
func (_m *DB) GetValues(dbName string, key string) (*types.GetHistoricalDataResponseEnvelope, error) {
	ret := _m.Called(dbName, key)

	var r0 *types.GetHistoricalDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetHistoricalDataResponseEnvelope); ok {
		r0 = rf(dbName, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetHistoricalDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dbName, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValuesDeletedByUser provides a mock function with given fields: userID
func (_m *DB) GetValuesDeletedByUser(userID string) (*types.GetDataProvenanceResponseEnvelope, error) {
	ret := _m.Called(userID)

	var r0 *types.GetDataProvenanceResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetDataProvenanceResponseEnvelope); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataProvenanceResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValuesReadByUser provides a mock function with given fields: userID
func (_m *DB) GetValuesReadByUser(userID string) (*types.GetDataProvenanceResponseEnvelope, error) {
	ret := _m.Called(userID)

	var r0 *types.GetDataProvenanceResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetDataProvenanceResponseEnvelope); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataProvenanceResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValuesWrittenByUser provides a mock function with given fields: userID
func (_m *DB) GetValuesWrittenByUser(userID string) (*types.GetDataProvenanceResponseEnvelope, error) {
	ret := _m.Called(userID)

	var r0 *types.GetDataProvenanceResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetDataProvenanceResponseEnvelope); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataProvenanceResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWriters provides a mock function with given fields: dbName, key
func (_m *DB) GetWriters(dbName string, key string) (*types.GetDataWritersResponseEnvelope, error) {
	ret := _m.Called(dbName, key)

	var r0 *types.GetDataWritersResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetDataWritersResponseEnvelope); ok {
		r0 = rf(dbName, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataWritersResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dbName, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Height provides a mock function with given fields:
func (_m *DB) Height() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDBExists provides a mock function with given fields: name
func (_m *DB) IsDBExists(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsLeader provides a mock function with given fields:
func (_m *DB) IsLeader() *errors.NotLeaderError {
	ret := _m.Called()

	var r0 *errors.NotLeaderError
	if rf, ok := ret.Get(0).(func() *errors.NotLeaderError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.NotLeaderError)
		}
	}

	return r0
}

// LedgerHeight provides a mock function with given fields:
func (_m *DB) LedgerHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitTransaction provides a mock function with given fields: tx, timeout
func (_m *DB) SubmitTransaction(tx interface{}, timeout time.Duration) (*types.TxReceiptResponseEnvelope, error) {
	ret := _m.Called(tx, timeout)

	var r0 *types.TxReceiptResponseEnvelope
	if rf, ok := ret.Get(0).(func(interface{}, time.Duration) *types.TxReceiptResponseEnvelope); ok {
		r0 = rf(tx, timeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.TxReceiptResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, time.Duration) error); ok {
		r1 = rf(tx, timeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
