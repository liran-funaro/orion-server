// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hyperledger-labs/orion-server/pkg/marshal"
	"github.com/hyperledger-labs/orion-server/pkg/types"
	"google.golang.org/protobuf/proto"
)

const MultiPartFormData = "multipart/form-data"

// SendHTTPResponse writes HTTP response back including HTTP code number and encode payload
func SendHTTPResponse(w http.ResponseWriter, code int, payload interface{}) {
	var response []byte
	if p, ok := payload.(proto.Message); ok {
		response, _ = marshal.DefaultMarshaler().Marshal(p)
	} else {
		response, _ = json.Marshal(payload)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Printf("Warning: failed to write response [%v] to the response writer\n", w)
	}
}

// SendHTTPRedirectServer replaces the Host in the request URL with hostPort, and redirects using
// StatusTemporaryRedirect (307).
func SendHTTPRedirectServer(w http.ResponseWriter, r *http.Request, hostPort string) {
	u, err := url.ParseRequestURI(r.URL.String())
	if err != nil {
		SendHTTPResponse(w, http.StatusInternalServerError,
			&types.HttpResponseErr{ErrMsg: fmt.Sprintf("cannot parse request URL: %s", err.Error())})
		return
	}

	u.Host = hostPort
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func GetStartAndEndBlockNum(params map[string]string) (uint64, uint64, error) {
	startBlockNum, err := GetUintParam("startId", params)
	if err != nil {
		return 0, 0, err
	}

	endBlockNum, err := GetUintParam("endId", params)
	if err != nil {
		return 0, 0, err
	}

	if endBlockNum < startBlockNum {
		return 0, 0, &types.HttpResponseErr{
			ErrMsg: fmt.Sprintf("query error: startId=%d > endId=%d", startBlockNum, endBlockNum),
		}
	}

	return startBlockNum, endBlockNum, nil
}

func GetUintParam(key string, params map[string]string) (uint64, *types.HttpResponseErr) {
	valStr, ok := params[key]
	if !ok {
		return 0, &types.HttpResponseErr{
			ErrMsg: "query error - bad or missing literal: " + key,
		}
	}
	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, &types.HttpResponseErr{
			ErrMsg: "query error - bad or missing literal: " + key + " " + err.Error(),
		}
	}
	return val, nil
}

func GetBlockNumAndTxIndex(params map[string]string) (uint64, uint64, error) {
	blockNum, err := GetUintParam("blockId", params)
	if err != nil {
		return 0, 0, err
	}

	txIndex, err := GetUintParam("idx", params)
	if err != nil {
		return 0, 0, err
	}

	return blockNum, txIndex, nil
}

func GetBlockNum(params map[string]string) (uint64, error) {
	blockNum, err := GetUintParam("blockId", params)
	if err != nil {
		return 0, err
	}

	return blockNum, nil
}

func GetVersion(params map[string]string) (*types.Version, error) {
	if _, ok := params["blknum"]; !ok {
		return nil, nil
	}

	blockNum, err := GetUintParam("blknum", params)
	if err != nil {
		return nil, err
	}

	txNum, err := GetUintParam("txnum", params)
	if err != nil {
		return nil, err
	}

	return &types.Version{
		BlockNum: blockNum,
		TxNum:    txNum,
	}, nil
}

func GetBase64urlKey(params map[string]string, name string) (string, error) {
	base64urlKey, ok := params[name]
	if !ok {
		return "", &types.HttpResponseErr{ErrMsg: fmt.Sprintf("Missing key: %s (in base64 URL encoding)", name)}
	}
	keyBytes, err := base64.RawURLEncoding.DecodeString(base64urlKey)
	if err != nil {
		return "", &types.HttpResponseErr{ErrMsg: fmt.Sprintf("Failed to decode base64 URL key: %s: %s", name, err.Error())}
	}

	return string(keyBytes), nil
}
