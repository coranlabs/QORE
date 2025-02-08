package sbi

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	// "runtime/debug"
	"sync"
	"time"

	logger_util "github.com/coranlabs/CORAN_LIB_LOGGER_UTIL"
	"github.com/coranlabs/CORAN_LIB_UTIL/httpwrapper"
	"github.com/coranlabs/CORAN_SCP/Application_entity/logger"
	"github.com/coranlabs/CORAN_SCP/Application_entity/app"
	"github.com/gin-gonic/gin"
)

type ServerScp interface {
	app.App

	// Consumer() *consumer.Consumer
	// Processor() *processor.Processor
}

type Server struct {
	ServerScp

	httpServer *http.Server
	router     *gin.Engine
}

func (s *Server) Run(wg *sync.WaitGroup) error {
	wg.Add(1)
	// log.Printf("yha takk aya")
	go s.startServer(wg)
	// log.Printf("yha takk aya 2")

	logger.SBILog.Infoln("SBI server started")
	return nil
}

func (s *Server) startServer(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.SBILog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			s.Terminate()
		}
		wg.Done()
	}()
	// TODO: support TLS mutual authentication for OAuth
	cfg := s.Config()
	serverScheme := cfg.GetSbiScheme()
	var err error // log.Printf("yha takk aya http")

	if serverScheme == "http" {
		s.httpServer.ListenAndServe()
	} else if serverScheme == "https" {
		err = s.httpServer.ListenAndServeTLS(cfg.GetScpCertPemPath(), cfg.GetScpPrivKeyPath())
	} else {
		err = fmt.Errorf("no support this scheme[%s]", serverScheme)
	}
	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}

// func newFunction(s *Server) {
// 	// log.Printf("yha takk aya server")

// 	cfg := s.Config()
// 	serverScheme := cfg.GetSbiScheme()
// 	// log.Printf("yha takk aya http")
// 	var err error
// 	if serverScheme == "http" {
// 		s.httpServer.ListenAndServe()
// 	} else if serverScheme == "https" {

// 		err = s.httpServer.ListenAndServeTLS(
// 			cfg.GetScpCertPemPath(),
// 			cfg.GetScpPrivKeyPath())
// 	} else {
// 		err = fmt.Errorf("No support this scheme[%s]", serverScheme)
// 	}

//		if err != nil && err != http.ErrServerClosed {
//			logger.SBILog.Errorf("SBI server error: %v", err)
//		}
//		logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
//	}
func (s *Server) Stop() {
	// server stop
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(toCtx); err != nil {
		logger.SBILog.Errorf("Could not close SBI server: %#v", err)
	}
}

func NewServer(scp ServerScp, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerScp: scp,
		router:    logger_util.NewGinWithLogrus(logger.GinLog),
	}
	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)

	s.applyService()

	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	return s, nil
}

func (s *Server) applyService() {
	accesstokenRoutes := s.getIndexroutes()
	accesstokenGroup := s.router.Group("") // accesstoken service didn't have api prefix
	applyRoutes(accesstokenGroup, accesstokenRoutes)

	// Run routes for AUSF in a separate goroutine
	s.router.Group(AusfAuthResUriPrefix).Any("/*action", s.handleNausfRoutes)
	s.router.Group(AusfSorprotectionResUriPrefix).Any("/*action", s.handleNausfRoutes)
	s.router.Group(AusfUpuprotectionResUriPrefix).Any("/*action", s.handleNausfRoutes)

	// Run routes for PCF in a separate goroutine
	s.router.Group(PcfAMpolicyCtlResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfPolicyAuthResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfCallbackResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfSMpolicyCtlResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfBdtPolicyCtlResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfOamResUriPrefix).Any("/*action", s.handleNpcfRoutes)
	s.router.Group(PcfUePolicyCtlResUriPrefix).Any("/*action", s.handleNpcfRoutes)

	// Run routes for UDM in a separate goroutine
	s.router.Group(UdmSorprotectionResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmAuthResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmfUpuprotectionResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmEcmResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmSdmResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmEeResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmUecmResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmPpResUriPrefix).Any("/*action", s.handleNudmRoutes)
	s.router.Group(UdmUeauResUriPrefix).Any("/*action", s.handleNudmRoutes)

	// Run routes for SMF in a separate goroutine
	s.router.Group(SmfEventExposureResUriPrefix).Any("/*action", s.handleNsmfRoutes)
	s.router.Group(SmfPdusessionResUriPrefix).Any("/*action", s.handleNsmfRoutes)
	s.router.Group(SmfOamUriPrefix).Any("/*action", s.handleNsmfRoutes)
	s.router.Group(SmfCallbackUriPrefix).Any("/*action", s.handleNsmfRoutes)

	// Run routes for NSSF in a separate goroutine
	s.router.Group(NssfNssaiavailResUriPrefix).Any("/*action", s.handleNnssfRoutes)
	s.router.Group(NssfNsselectResUriPrefix).Any("/*action", s.handleNnssfRoutes)

	// Run routes for AMF in a separate goroutine
	s.router.Group(AmfCallbackResUriPrefix).Any("/*action", s.handleNamfRoutes)
	s.router.Group(AmfCommResUriPrefix).Any("/*action", s.handleNamfRoutes)
	s.router.Group(AmfEvtsResUriPrefix).Any("/*action", s.handleNamfRoutes)
	s.router.Group(AmfLocResUriPrefix).Any("/*action", s.handleNamfRoutes)
	s.router.Group(AmfMtResUriPrefix).Any("/*action", s.handleNamfRoutes)
	s.router.Group(AmfOamResUriPrefix).Any("/*action", s.handleNamfRoutes)

	s.router.Group(UdrDrResUriPrefix).Any("/*action", s.handleNudrRoutes)
	s.router.Group(ConvergedChargingResUriPrefix).Any("/*action", s.handleChfRoutes)
}
