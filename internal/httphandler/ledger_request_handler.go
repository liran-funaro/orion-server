package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.ibm.com/blockchaindb/server/internal/bcdb"
	"github.ibm.com/blockchaindb/server/pkg/constants"
	"github.ibm.com/blockchaindb/server/pkg/cryptoservice"
	"github.ibm.com/blockchaindb/server/pkg/logger"
)

// ledgerRequestHandler handles query associated with the
// chain of blocks
type ledgerRequestHandler struct {
	db          backend.DB
	sigVerifier *cryptoservice.SignatureVerifier
	router      *mux.Router
	logger      *logger.SugarLogger
}

// NewLedgerRequestHandler creates users request handler
func NewLedgerRequestHandler(db backend.DB, logger *logger.SugarLogger) http.Handler {
	handler := &ledgerRequestHandler{
		db:          db,
		sigVerifier: cryptoservice.NewVerifier(db),
		router:      mux.NewRouter(),
		logger:      logger,
	}

	// HTTP GET "/ledger/block/{blockId}" gets block header
	handler.router.HandleFunc(constants.GetBlockHeader, handler.blockQuery).Methods(http.MethodGet)
	// HTTP GET "/ledger/path?start={startId}&end={endId}" gets shortest path between blocks
	handler.router.HandleFunc(constants.GetPath, handler.pathQuery).Methods(http.MethodGet).Queries("start", "{startId:[0-9]+}", "end", "{endId:[0-9]+}")
	// HTTP GET "/ledger/proof/{blockId}?idx={idx}" gets proof for tx with idx index inside block blockId
	handler.router.HandleFunc(constants.GetTxProof, handler.txProof).Methods(http.MethodGet).Queries("idx", "{idx:[0-9]+}")
	// HTTP GET "/ledger/tx/receipt/{txId}" gets transaction receipt
	handler.router.HandleFunc(constants.GetTxReceipt, handler.txReceipt).Methods(http.MethodGet)
	// HTTP GET "/ledger/path?start={startId}&end={endId}" with invalid query params
	handler.router.HandleFunc(constants.GetPath, handler.invalidPathQuery).Methods(http.MethodGet)
	// HTTP GET "/ledger/proof/{blockId}?idx={idx}" with invalid query params
	handler.router.HandleFunc(constants.GetTxProof, handler.invalidTxProof).Methods(http.MethodGet)

	return handler
}

func (p *ledgerRequestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	p.router.ServeHTTP(responseWriter, request)
}

func (p *ledgerRequestHandler) blockQuery(response http.ResponseWriter, request *http.Request) {
	queryEnv, respondedErr := extractBlockQueryEnvelope(request, response)
	if respondedErr {
		return
	}

	err, code := VerifyRequestSignature(p.sigVerifier, queryEnv.Payload.UserID, queryEnv.Signature, queryEnv.Payload)
	if err != nil {
		SendHTTPResponse(response, code, err)
		return
	}

	data, err := p.db.GetBlockHeader(queryEnv.Payload.UserID, queryEnv.Payload.BlockNumber)
	if err != nil {
		var status int

		switch err.(type) {
		case *backend.PermissionErr:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError // TODO deal with 404 not found, it's not a 5xx
		}

		SendHTTPResponse(
			response,
			status,
			&ResponseErr{
				ErrMsg: "error while processing '" + request.Method + " " + request.URL.String() + "' because " + err.Error(),
			})
		return
	}

	SendHTTPResponse(response, http.StatusOK, data)
}

func (p *ledgerRequestHandler) pathQuery(response http.ResponseWriter, request *http.Request) {
	queryEnv, respondedErr := extractLedgerPathQueryEnvelope(request, response)
	if respondedErr {
		return
	}

	err, code := VerifyRequestSignature(p.sigVerifier, queryEnv.Payload.UserID, queryEnv.Signature, queryEnv.Payload)
	if err != nil {
		SendHTTPResponse(response, code, err)
		return
	}

	data, err := p.db.GetLedgerPath(queryEnv.Payload.UserID, queryEnv.Payload.StartBlockNumber, queryEnv.Payload.EndBlockNumber)
	if err != nil {
		var status int

		switch err.(type) {
		case *backend.PermissionErr:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
		}

		SendHTTPResponse(
			response,
			status,
			&ResponseErr{
				ErrMsg: "error while processing '" + request.Method + " " + request.URL.String() + "' because " + err.Error(),
			})
		return
	}

	SendHTTPResponse(response, http.StatusOK, data)
}

func (p *ledgerRequestHandler) txProof(response http.ResponseWriter, request *http.Request) {
	queryEnv, respondedErr := extractGetTxProofQueryEnvelope(request, response)
	if respondedErr {
		return
	}

	err, code := VerifyRequestSignature(p.sigVerifier, queryEnv.Payload.UserID, queryEnv.Signature, queryEnv.Payload)
	if err != nil {
		SendHTTPResponse(response, code, err)
		return
	}

	data, err := p.db.GetTxProof(queryEnv.GetPayload().GetUserID(), queryEnv.GetPayload().GetBlockNumber(), queryEnv.GetPayload().GetTxIndex())
	if err != nil {
		var status int

		switch err.(type) {
		case *backend.PermissionErr:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError // TODO deal with 404 not found, it's not a 5xx
		}

		SendHTTPResponse(
			response,
			status,
			&ResponseErr{
				ErrMsg: "error while processing '" + request.Method + " " + request.URL.String() + "' because " + err.Error(),
			})
		return
	}

	SendHTTPResponse(response, http.StatusOK, data)
}

func (p *ledgerRequestHandler) txReceipt(response http.ResponseWriter, request *http.Request) {
	queryEnv, respondedErr := extractGetTxReceiptQueryEnvelope(request, response)
	if respondedErr {
		return
	}

	err, code := VerifyRequestSignature(p.sigVerifier, queryEnv.Payload.UserID, queryEnv.Signature, queryEnv.Payload)
	if err != nil {
		SendHTTPResponse(response, code, err)
		return
	}

	data, err := p.db.GetTxReceipt(queryEnv.GetPayload().GetUserID(), queryEnv.GetPayload().GetTxID())
	if err != nil {
		var status int

		switch err.(type) {
		case *backend.PermissionErr:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError // TODO deal with 404 not found, it's not a 5xx
		}

		SendHTTPResponse(
			response,
			status,
			&ResponseErr{
				ErrMsg: "error while processing '" + request.Method + " " + request.URL.String() + "' because " + err.Error(),
			})
		return
	}

	SendHTTPResponse(response, http.StatusOK, data)
}

func (p *ledgerRequestHandler) invalidPathQuery(response http.ResponseWriter, request *http.Request) {
	err := &ResponseErr{
		ErrMsg: "query error - bad or missing start/end block number",
	}
	SendHTTPResponse(response, http.StatusBadRequest, err)
}

func (p *ledgerRequestHandler) invalidTxProof(response http.ResponseWriter, request *http.Request) {
	err := &ResponseErr{
		ErrMsg: "query error - bad or missing tx index",
	}
	SendHTTPResponse(response, http.StatusBadRequest, err)
}

func getUintParam(key string, params map[string]string) (uint64, *ResponseErr) {
	valStr, ok := params[key]
	if !ok {
		return 0, &ResponseErr{
			ErrMsg: "query error - bad or missing block number literal" + key,
		}
	}
	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, &ResponseErr{
			ErrMsg: "query error - bad or missing block number literal " + key + " " + err.Error(),
		}
	}
	return val, nil
}