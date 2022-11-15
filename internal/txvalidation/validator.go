// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package txvalidation

import (
	"sync"
	"time"

	"github.com/hyperledger-labs/orion-server/internal/identity"
	"github.com/hyperledger-labs/orion-server/internal/utils"
	"github.com/hyperledger-labs/orion-server/internal/worldstate"
	"github.com/hyperledger-labs/orion-server/pkg/cryptoservice"
	"github.com/hyperledger-labs/orion-server/pkg/logger"
	"github.com/hyperledger-labs/orion-server/pkg/marshal"
	"github.com/hyperledger-labs/orion-server/pkg/types"
	"github.com/pkg/errors"
)

// Validator validates the each transaction read set present in a
// block against the committed version to ensure the requested
// isolation level
type Validator struct {
	configTxValidator    *ConfigTxValidator
	dbAdminTxValidator   *dbAdminTxValidator
	userAdminTxValidator *userAdminTxValidator
	dataTxValidator      *dataTxValidator
	signValidator        *txSigValidator
	logger               *logger.SugarLogger
}

type Config struct {
	DB     worldstate.DB
	Logger *logger.SugarLogger
}

// NewValidator creates a new Validator
func NewValidator(conf *Config, disableCommitPhaseValidation bool) *Validator {
	idQuerier := identity.NewQuerier(conf.DB)
	txSigValidator := &txSigValidator{
		sigVerifier: cryptoservice.NewVerifier(idQuerier, conf.Logger),
		logger:      conf.Logger,
		marshaler:   marshal.DefaultMarshaler(),
	}

	return &Validator{
		configTxValidator: &ConfigTxValidator{
			db:              conf.DB,
			identityQuerier: idQuerier,
			sigValidator:    txSigValidator,
			logger:          conf.Logger,
		},

		dbAdminTxValidator: &dbAdminTxValidator{
			db:              conf.DB,
			identityQuerier: idQuerier,
			sigValidator:    txSigValidator,
			logger:          conf.Logger,
		},

		userAdminTxValidator: &userAdminTxValidator{
			db:              conf.DB,
			identityQuerier: idQuerier,
			sigValidator:    txSigValidator,
			logger:          conf.Logger,
		},

		dataTxValidator: &dataTxValidator{
			db:                           conf.DB,
			identityQuerier:              idQuerier,
			sigValidator:                 txSigValidator,
			logger:                       conf.Logger,
			disableCommitPhaseValidation: disableCommitPhaseValidation,
		},

		signValidator: txSigValidator,

		logger: conf.Logger,
	}
}

// ValidateBlock validates each transaction present in the block to ensure
// the request isolation level
func (v *Validator) ValidateBlock(block *types.Block) ([]*types.ValidationInfo, error) {
	if block.Header.BaseHeader.Number == 1 {
		// for the genesis block, which is created by the node itself, we cannot
		// do a regular validation, but we still need to validate the entries.
		return v.configTxValidator.validateGenesis(block.GetConfigTxEnvelope())
	}

	switch block.Payload.(type) {
	case *types.Block_DataTxEnvelopes:
		dataTxEnvs := block.GetDataTxEnvelopes().Envelopes
		start := time.Now()
		valInfoArray, usersWithValidSigPerTX, err := v.parallelSigValidation(dataTxEnvs)
		utils.Stats.UpdateSigValidationTime(time.Since(start))
		if err != nil {
			return nil, err
		}

		pendingOps := newPendingOperations()
		for txNum, txEnv := range dataTxEnvs {
			if valInfoArray[txNum].Flag != types.Flag_VALID {
				continue
			}

			valRes, err := v.dataTxValidator.validate(txEnv, usersWithValidSigPerTX[txNum], pendingOps)
			if err != nil {
				return nil, errors.WithMessage(err, "error while validating data transaction")
			}

			valInfoArray[txNum] = valRes
			if valRes.Flag != types.Flag_VALID {
				v.logger.Debugf("data transaction [%v] is invalid due to [%s]", txEnv.Payload, valRes.ReasonIfInvalid)
				continue
			}

			for _, ops := range txEnv.Payload.DbOperations {
				for _, w := range ops.DataWrites {
					pendingOps.addWrite(ops.DbName, w.Key)
				}

				for _, d := range ops.DataDeletes {
					pendingOps.addDelete(ops.DbName, d.Key)
				}
			}
		}

		return valInfoArray, nil

	case *types.Block_UserAdministrationTxEnvelope:
		userTxEnv := block.GetUserAdministrationTxEnvelope()
		valRes, err := v.userAdminTxValidator.validate(userTxEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "error while validating user administrative transaction")
		}

		if valRes.Flag != types.Flag_VALID {
			v.logger.Debugf("user administration transaction [%v] is invalid due to [%s]", userTxEnv.Payload, valRes.ReasonIfInvalid)
		}

		return []*types.ValidationInfo{
			valRes,
		}, nil

	case *types.Block_DbAdministrationTxEnvelope:
		dbTxEnv := block.GetDbAdministrationTxEnvelope()
		valRes, err := v.dbAdminTxValidator.validate(dbTxEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "error while validating db administrative transaction")
		}

		if valRes.Flag != types.Flag_VALID {
			v.logger.Debugf("database administration transaction [%v] is invalid due to [%s]", dbTxEnv.Payload, valRes.ReasonIfInvalid)
		}

		return []*types.ValidationInfo{
			valRes,
		}, nil

	case *types.Block_ConfigTxEnvelope:
		configTxEnv := block.GetConfigTxEnvelope()
		valRes, err := v.configTxValidator.Validate(configTxEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "error while validating config transaction")
		}

		if valRes.Flag != types.Flag_VALID {
			v.logger.Debugf("cluster config transaction [%v] is invalid due to [%s]", configTxEnv, valRes.ReasonIfInvalid)
		}

		return []*types.ValidationInfo{
			valRes,
		}, nil

	default:
		return nil, errors.Errorf("unexpected transaction envelope in the block")
	}
}

// ConfigValidator provides a pointer to the internal validator that verifies config transactions.
func (v *Validator) ConfigValidator() *ConfigTxValidator {
	return v.configTxValidator
}

func (v *Validator) parallelSigValidation(dataTxEnvs []*types.DataTxEnvelope) ([]*types.ValidationInfo, [][]string, error) {
	valInfoPerTx := make([]*types.ValidationInfo, len(dataTxEnvs))
	usersWithValidSigPerTX := make([][]string, len(dataTxEnvs))
	errorPerTx := make([]error, len(dataTxEnvs))

	var wg sync.WaitGroup
	wg.Add(len(dataTxEnvs))

	for txNumber, txEnvelope := range dataTxEnvs {
		go func(txEnv *types.DataTxEnvelope, txNum int) {
			defer wg.Done()

			usersWithValidSignTx, vInfo, vErr := v.dataTxValidator.commitPhaseValidateSignatures(txEnv)

			if vErr != nil {
				errorPerTx[txNum] = vErr
				return
			}

			usersWithValidSigPerTX[txNum] = usersWithValidSignTx
			valInfoPerTx[txNum] = vInfo
			if vInfo.Flag != types.Flag_VALID {
				v.logger.Debugf("data transaction [%v] is invalid due to [%s]", txEnv.Payload, vInfo.ReasonIfInvalid)
			}
		}(txEnvelope, txNumber)
	}
	wg.Wait()

	for txNum, err := range errorPerTx {
		if err != nil {
			v.logger.Errorf("error validating signatures in tx number %d, error: %s", txNum, err)
			return nil, nil, err
		}
	}
	return valInfoPerTx, usersWithValidSigPerTX, nil
}

type pendingOperations struct {
	pendingWrites  map[string]bool
	pendingDeletes map[string]bool
}

func newPendingOperations() *pendingOperations {
	return &pendingOperations{
		pendingWrites:  make(map[string]bool),
		pendingDeletes: make(map[string]bool),
	}
}

func (p *pendingOperations) addWrite(dbName, key string) {
	ckey := constructCompositeKey(dbName, key)
	p.pendingWrites[ckey] = true
}

func (p *pendingOperations) addDelete(dbName, key string) {
	ckey := constructCompositeKey(dbName, key)
	p.pendingDeletes[ckey] = true
}

func (p *pendingOperations) existDelete(dbName, key string) bool {
	ckey := constructCompositeKey(dbName, key)
	return p.pendingDeletes[ckey]
}

func (p *pendingOperations) exist(dbName, key string) bool {
	ckey := constructCompositeKey(dbName, key)
	return p.pendingWrites[ckey] || p.pendingDeletes[ckey]
}

func constructCompositeKey(dbName, key string) string {
	return dbName + "~" + key
}
