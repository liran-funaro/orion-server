// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"sync"

	"github.com/hyperledger-labs/orion-server/pkg/logger"
	"github.com/hyperledger-labs/orion-server/pkg/types"
)

type PendingTxs struct {
	//sync.RWMutex
	txs sync.Map
	//txs map[string]*CompletionPromise

	logger *logger.SugarLogger
}

func NewPendingTxs(logger *logger.SugarLogger) *PendingTxs {
	return &PendingTxs{
		//txs:    make(map[string]*CompletionPromise),
		logger: logger,
	}
}

func (p *PendingTxs) Add(txID string, promise *CompletionPromise) bool {
	//p.Lock()
	//defer p.Unlock()

	_, loaded := p.txs.LoadOrStore(txID, promise)
	p.txs.Store(txID, promise)
	//p.txs[txID] = promise
	return loaded
}

func (p *PendingTxs) DeleteWithNoAction(txID string) {
	//p.Lock()
	//defer p.Unlock()

	p.txs.Delete(txID)
	//p.txs[txID] = promise
}

func (p *PendingTxs) loadAndDelete(txID string) (*CompletionPromise, bool) {
	promise, loaded := p.txs.LoadAndDelete(txID)
	if !loaded {
		return nil, loaded
	}
	return promise.(*CompletionPromise), true
}

// DoneWithReceipt is called after the commit of a block.
// The `txIDs` slice must be in the same order that transactions appear in the block.
func (p *PendingTxs) DoneWithReceipt(txIDs []string, blockHeader *types.BlockHeader) {
	p.logger.Debugf("Done with receipt, block number: %d; txIDs: %v", blockHeader.GetBaseHeader().GetNumber(), txIDs)

	//p.Lock()
	//defer p.Unlock()

	for txIndex, txID := range txIDs {
		promise, loaded := p.loadAndDelete(txID)
		if !loaded {
			continue
		}
		promise.done(
			&types.TxReceipt{
				Header:  blockHeader,
				TxIndex: uint64(txIndex),
			},
		)
	}
}

// ReleaseWithError is called when block replication fails with an error, typically NotLeaderError.
// This may come from the block replicator or the block creator.
// The `txIDs` slice does not have to be in the same order that transactions appear in the block.
func (p *PendingTxs) ReleaseWithError(txIDs []string, err error) {
	p.logger.Debugf("Release with error: %s; txIDs: %v", err, txIDs)

	//p.Lock()
	//defer p.Unlock()

	for _, txID := range txIDs {
		promise, loaded := p.loadAndDelete(txID)
		if !loaded {
			continue
		}
		promise.error(err)
	}
}

func (p *PendingTxs) Has(txID string) bool {
	//p.RLock()
	//defer p.RUnlock()

	_, ok := p.txs.Load(txID)
	return ok
}

func (p *PendingTxs) Empty() bool {
	//p.RLock()
	//defer p.RUnlock()

	empty := true
	p.txs.Range(func(key, value interface{}) bool {
		empty = false
		return false
	})
	return empty
}
