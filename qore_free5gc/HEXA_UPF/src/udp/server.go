// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package udp

import (
	"context"
	"net/http"
	"time"

	"github.com/coranlabs/HEXA_UPF/src/logger"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	logger.AppLog.Infof("running on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
