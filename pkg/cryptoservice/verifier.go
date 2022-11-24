// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cryptoservice

import (
	"crypto/x509"
	"sync"

	"github.com/hyperledger-labs/orion-server/pkg/logger"

	"github.com/hyperledger-labs/orion-server/pkg/crypto"
)

//go:generate counterfeiter -o mocks/user_db_querier.go --fake-name UserDBQuerier . UserDBQuerier

type UserDBQuerier interface {
	GetCertificate(userID string) (*x509.Certificate, error)
}

func NewVerifier(userQuerier UserDBQuerier, logger *logger.SugarLogger) *SignatureVerifier {
	return &SignatureVerifier{
		userDBQuerier: userQuerier,
		logger:        logger,
	}
}

//TODO keep a cache of user and parsed certificates to avoid going to the DB and parsing the certificate
// on every TX. Provide a mechanism to invalidate the cache when the user database changes.

type SignatureVerifier struct {
	userDBQuerier UserDBQuerier
	logger        *logger.SugarLogger
	certCache     sync.Map
}

func (sv *SignatureVerifier) GetCert(userID string) (*x509.Certificate, error) {
	//cert, ok := sv.certCache.Load(userID)
	//if ok {
	//	return cert.(*x509.Certificate), nil
	//}

	newCert, err := sv.userDBQuerier.GetCertificate(userID)
	if err != nil {
		return nil, err
	}
	//sv.certCache.Store(userID, newCert)
	return newCert, nil
}

func (sv *SignatureVerifier) Verify(userID string, signature, body []byte) error {
	cert, err := sv.GetCert(userID)
	if err != nil {
		sv.logger.Infof("Error during GetCertificate: userID: %s, error: %s", userID, err)
		return err
	}
	verifier := crypto.Verifier{Certificate: cert}
	if err = verifier.Verify(body, signature); err != nil {
		sv.logger.Infof("Failed to verify signature: userID: %s, error: %s", userID, err)
		return err
	}
	return err
}
