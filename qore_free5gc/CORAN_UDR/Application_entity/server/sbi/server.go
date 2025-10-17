package sbi

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
	processor "github.com/coranlabs/CORAN_UDR/Application_entity/server/sbi/processor"

	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	"github.com/coranlabs/CORAN_UDR/Application_entity/server/util"
	"github.com/lakshya-chopra/httpwrapper"
)

type Server struct {
	UDR

	httpServer *http.Server
	router     *gin.Engine
}

type UDR interface {
	app.App

	Processor() *processor.Processor
}

func NewServer(udr UDR, tlsKeyLogPath string) *Server {
	s := &Server{
		UDR: udr,
	}

	s.router = newRouter(s)

	server_cert, err1 := tls.LoadX509KeyPair(factory.UdrConfig.GetCertPemPath(), factory.UdrConfig.GetCertKeyPath())

	if err1 != nil {
		log.Fatal(err1)
	}
	server, err := bindRouter(udr, s.router, tlsKeyLogPath, server_cert)
	s.httpServer = server

	if err != nil {
		logger.SBILog.Errorf("bind Router Error: %+v", err)
		panic("Server initialization failed")
	}

	return s
}

func (s *Server) Run(wg *sync.WaitGroup) {
	logger.SBILog.Info("Starting server...")

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := s.serve()
		if err != http.ErrServerClosed {
			logger.SBILog.Panicf("HTTP server setup failed: %+v", err)
		}
		logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
	}()
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
	const shutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.SBILog.Errorf("HTTP server shutdown failed: %+v", err)
	}
}

func bindRouter(udr app.App, router *gin.Engine, tlsKeyLogPath string, cert tls.Certificate) (*http.Server, error) {
	sbiConfig := udr.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)

	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router, cert)
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	dataRepositoryGroup := router.Group(factory.UdrDrResUriPrefix)
	dataRepositoryGroup.Use(func(c *gin.Context) {
		util.NewRouterAuthorizationCheck(models.ServiceName_NUDR_DR).Check(c, s.Context())
	})
	dataRepositoryRoutes := s.getDataRepositoryRoutes()
	AddService(dataRepositoryGroup, dataRepositoryRoutes)
	return router
}

func (s *Server) unsecureServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) secureServe() error {
	sbiConfig := s.UDR.Config()

	pemPath := sbiConfig.GetCertPemPath()
	if pemPath == "" {
		pemPath = factory.UdrDefaultCertPemPath
	}

	keyPath := sbiConfig.GetCertKeyPath()
	if keyPath == "" {
		keyPath = factory.UdrDefaultPrivateKeyPath
	}

	return s.httpServer.ListenAndServeTLS(pemPath, keyPath)
}

func (s *Server) serve() error {
	sbiConfig := s.UDR.Config().Configuration.Sbi

	switch sbiConfig.Scheme {
	case "http":
		return s.unsecureServe()
	case "https":
		return s.secureServe()
	default:
		return fmt.Errorf("invalid SBI scheme: %s", sbiConfig.Scheme)
	}
}
