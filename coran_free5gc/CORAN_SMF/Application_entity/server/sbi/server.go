package sbi

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/consumer"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/processor"
	util_oauth "github.com/coranlabs/CORAN_SMF/Application_entity/util/oauth"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
	"github.com/lakshya-chopra/httpwrapper"
)

const (
	APPLICATION_JSON  = "application/json"
	MULTIPART_RELATED = "multipart/related"
)

type ServerSmf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerSmf

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(smf ServerSmf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerSmf: smf,
	}

	smf_context.InitSmfContext(factory.SmfConfig)
	// allocate id for each upf
	smf_context.AllocateUPFID()
	smf_context.InitSMFUERouting(factory.UERoutingConfig)

	s.router = newRouter(s)

	server_cert, err1 := tls.LoadX509KeyPair(factory.SmfConfig.GetCertPemPath(), factory.SmfConfig.GetCertKeyPath())

	if err1 != nil {
		log.Fatal(err1)
	}

	bindAddr := fmt.Sprintf("%s:%d", s.Context().BindingIPv4, s.Context().SBIPort)
	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router, server_cert); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}

	return s, nil
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	smfCallbackGroup := router.Group(factory.SmfCallbackUriPrefix)
	smfCallbackRoutes := s.getCallbackRoutes()
	applyRoutes(smfCallbackGroup, smfCallbackRoutes)

	upiGroup := router.Group(factory.UpiUriPrefix)
	upiRoutes := s.getUPIRoutes()
	applyRoutes(upiGroup, upiRoutes)

	for _, serviceName := range factory.SmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serviceName) {
		case models.ServiceName_NSMF_PDUSESSION:
			smfPDUSessionGroup := router.Group(factory.SmfPdusessionResUriPrefix)
			smfPDUSessionRoutes := s.getPDUSessionRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NSMF_PDUSESSION)
			smfPDUSessionGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, smf_context.GetSelf())
			})
			applyRoutes(smfPDUSessionGroup, smfPDUSessionRoutes)
		case models.ServiceName_NSMF_EVENT_EXPOSURE:
			smfEventExposureGroup := router.Group(factory.SmfEventExposureResUriPrefix)
			smfEventExposureRoutes := s.getEventExposureRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NSMF_EVENT_EXPOSURE)
			smfEventExposureGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, smf_context.GetSelf())
			})
			applyRoutes(smfEventExposureGroup, smfEventExposureRoutes)
		case models.ServiceName_NSMF_OAM:
			smfOAMGroup := router.Group(factory.SmfOamUriPrefix)
			smfOAMRoutes := s.getOAMRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NSMF_OAM)
			smfOAMGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, smf_context.GetSelf())
			})
			applyRoutes(smfOAMGroup, smfOAMRoutes)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	err := s.Consumer().RegisterNFInstance(context.Background())
	if err != nil {
		retry_err := s.Consumer().RetrySendNFRegistration(10)
		if retry_err != nil {
			logger.InitLog.Errorln(retry_err)
			return err
		}
	}

	wg.Add(1)
	go s.startServer(wg)

	return nil
}

func (s *Server) Stop() {
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer != nil {
		logger.SBILog.Infof("Stop SBI server (listen on %s)", s.httpServer.Addr)
		toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(toCtx); err != nil {
			logger.SBILog.Errorf("Could not close SBI server: %#v", err)
		}
	}
}

func (s *Server) startServer(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.SBILog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		wg.Done()
	}()

	logger.SBILog.Infof("Start SBI server (listen on %s)", s.httpServer.Addr)

	var err error
	cfg := s.Config()
	scheme := cfg.GetSbiScheme()
	if scheme == "http" {
		err = s.httpServer.ListenAndServe()
	} else if scheme == "https" {
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetCertPemPath(),
			cfg.GetCertKeyPath())
	} else {
		err = fmt.Errorf("no support this scheme[%s]", scheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}
