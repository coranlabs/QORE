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
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/logger"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/server/sbi/processor"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/util"
	"github.com/lakshya-chopra/httpwrapper"
)

type nssfApp interface {
	app.NssfApp

	Processor() *processor.Processor
}

type Server struct {
	nssfApp

	httpServer *http.Server
	router     *gin.Engine
	processor  *processor.Processor
}

func NewServer(nssf nssfApp, tlsKeyLogPath string) *Server {
	s := &Server{
		nssfApp:   nssf,
		processor: nssf.Processor(),
	}

	s.router = newRouter(s)

	server_cert, err1 := tls.LoadX509KeyPair(factory.NssfDefaultCertPemPath, factory.NssfDefaultPrivateKeyPath)

	if err1 != nil {
		log.Fatal(err1)
	}

	server, err := bindRouter(nssf, s.router, tlsKeyLogPath, server_cert)
	s.httpServer = server

	if err != nil {
		logger.SBILog.Errorf("bind Router Error: %+v", err)
		panic("Server initialization failed")
	}

	return s
}

func (s *Server) Processor() *processor.Processor {
	return s.processor
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

func bindRouter(nssf app.NssfApp, router *gin.Engine, tlsKeyLogPath string, cert tls.Certificate) (*http.Server, error) {

	sbiConfig := nssf.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)

	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router, cert)
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	for _, serviceName := range s.Config().Configuration.ServiceNameList {
		switch serviceName {
		case models.ServiceName_NNSSF_NSSAIAVAILABILITY:
			nssaiAvailabilityGroup := router.Group(factory.NssfNssaiavailResUriPrefix)
			nssaiAvailabilityGroup.Use(func(c *gin.Context) {
				// oauth middleware
				util.NewRouterAuthorizationCheck(serviceName).Check(c, s.Context())
			})
			nssaiAvailabilityRoutes := s.getNssaiAvailabilityRoutes()
			AddService(nssaiAvailabilityGroup, nssaiAvailabilityRoutes)

		case models.ServiceName_NNSSF_NSSELECTION:
			nsSelectionGroup := router.Group(factory.NssfNsselectResUriPrefix)
			nsSelectionGroup.Use(func(c *gin.Context) {
				// oauth middleware
				util.NewRouterAuthorizationCheck(serviceName).Check(c, s.Context())
			})
			nsSelectionRoutes := s.getNsSelectionRoutes()
			AddService(nsSelectionGroup, nsSelectionRoutes)

		default:
			logger.SBILog.Warnf("Unsupported service name: %s", serviceName)
		}
	}

	return router
}

func (s *Server) unsecureServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) secureServe() error {
	sbiConfig := s.Config().Configuration.Sbi

	pemPath := sbiConfig.Tls.Pem
	if pemPath == "" {
		pemPath = factory.NssfDefaultCertPemPath
	}

	keyPath := sbiConfig.Tls.Key
	if keyPath == "" {
		keyPath = factory.NssfDefaultPrivateKeyPath
	}

	return s.httpServer.ListenAndServeTLS(pemPath, keyPath)
}

func (s *Server) serve() error {
	sbiConfig := s.Config().Configuration.Sbi

	switch sbiConfig.Scheme {
	case "http":
		return s.unsecureServe()
	case "https":
		return s.secureServe()
	default:
		return fmt.Errorf("invalid SBI scheme: %s", sbiConfig.Scheme)
	}
}
