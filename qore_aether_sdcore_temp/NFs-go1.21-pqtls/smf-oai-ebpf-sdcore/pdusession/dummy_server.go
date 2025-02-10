// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package pdusession

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/lakshya-chopra/http2_util"
	"github.com/omec-project/logger_util"
	"github.com/omec-project/path_util"
	"github.com/omec-project/smf/logger"
	"github.com/omec-project/smf/pfcp"
	"github.com/omec-project/smf/pfcp/udp"
)

func DummyServer() {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	AddService(router)

	go udp.Run(pfcp.Dispatch)

	smfKeyLogPath := path_util.Free5gcPath("free5gc/smfsslkey.log")
	smfPemPath := path_util.Free5gcPath("free5gc/support/TLS/smf.pem")
	smfkeyPath := path_util.Free5gcPath("free5gc/support/TLS/key.pem")

	var server *http.Server

	server_cert, err := tls.LoadX509KeyPair(smfPemPath, smfkeyPath)
	if err != nil {
		log.Fatal(err)
	}
	if srv, err := http2_util.NewServer(":29502", smfKeyLogPath, router, server_cert); err != nil {
	} else {
		server = srv
	}

	if err := server.ListenAndServeTLS(smfPemPath, smfkeyPath); err != nil {
		log.Fatal(err)
	}
}
