package sbi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/coranlabs/CORAN_AMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_AMF/Application_entity/server/sbi/consumer"
	"github.com/coranlabs/CORAN_AMF/Application_entity/server/sbi/processor"
	amf_context "github.com/coranlabs/CORAN_AMF/Messages_controller/context"

	"github.com/coranlabs/CORAN_AMF/Application_entity/config/factory"
	"github.com/coranlabs/CORAN_AMF/Application_entity/pkg/app"
	util_oauth "github.com/coranlabs/CORAN_AMF/Application_entity/server/util"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/lakshya-chopra/httpwrapper"
)

var (
	reqbody         = "[Request Body] "
	applicationjson = "application/json"
	multipartrelate = "multipart/related"
)

type ServerAmf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerAmf

	httpServer *http.Server
	router     *gin.Engine
}

func PrintCertificateDetails(cert *x509.Certificate) {

	sep := strings.Repeat("-", 15)

	fmt.Printf("\n%s Server Certificate%s\n", sep, sep)

	fmt.Printf("Subject: %s\n", cert.Subject)
	fmt.Printf("Issuer: %s\n", cert.Issuer)
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber)
	fmt.Printf("Not Before: %s\n", cert.NotBefore)
	fmt.Printf("Not After: %s\n", cert.NotAfter)
	fmt.Printf("Key Usage: %x\n", cert.KeyUsage)
	fmt.Printf("Ext Key Usage: %v\n", cert.ExtKeyUsage)
	fmt.Printf("DNS Names: %v\n", cert.DNSNames)
	// fmt.Printf("Email Addresses: %v\n", cert.EmailAddresses)
	fmt.Printf("IP Addresses: %v\n", cert.IPAddresses)
	// fmt.Printf("URIs: %v\n", cert.URIs)
	fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm)

	fmt.Println("\nPEM Encoded Certificate:")
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	fmt.Println(string(pemBytes))

	fmt.Printf("%s End %s", sep, sep)
}

func ReadCertificate(filename string) (*x509.Certificate, error) {
	// Read the certificate file
	certPEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

func NewServer(amf ServerAmf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerAmf: amf,
	}

	s.router = newRouter(s)

	cfg := s.Config()

	server_cert, err1 := tls.LoadX509KeyPair(cfg.GetCertPemPath(), cfg.GetCertKeyPath())

	if err1 != nil {
		log.Fatal(err1)
	}

	cert, err2 := ReadCertificate(cfg.GetCertPemPath())
	if err2 != nil {
		log.Fatal(err2)
	} else {
		PrintCertificateDetails(cert)
	}

	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)
	var err error
	if s.httpServer, err = httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, s.router, server_cert); err != nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}
	s.httpServer.ErrorLog = log.New(logger.SBILog.WriterLevel(logrus.ErrorLevel), "HTTP2: ", 0)

	return s, nil
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	amfHttpCallBackGroup := router.Group(factory.AmfCallbackResUriPrefix)
	amfHttpCallBackRoutes := s.getHttpCallBackRoutes()
	applyRoutes(amfHttpCallBackGroup, amfHttpCallBackRoutes)

	for _, serverName := range factory.AmfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serverName) {
		case models.ServiceName_NAMF_COMM:
			amfCommunicationGroup := router.Group(factory.AmfCommResUriPrefix)
			amfCommunicationRoutes := s.getCommunicationRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_COMM)
			amfCommunicationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfCommunicationGroup, amfCommunicationRoutes)
		case models.ServiceName_NAMF_EVTS:
			amfEventExposureGroup := router.Group(factory.AmfEvtsResUriPrefix)
			amfEventExposureRoutes := s.getEventexposureRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_EVTS)
			amfEventExposureGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfEventExposureGroup, amfEventExposureRoutes)
		case models.ServiceName_NAMF_MT:
			amfMTGroup := router.Group(factory.AmfMtResUriPrefix)
			amfMTRoutes := s.getMTRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_MT)
			amfMTGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfMTGroup, amfMTRoutes)
		case models.ServiceName_NAMF_LOC:
			amfLocationGroup := router.Group(factory.AmfLocResUriPrefix)
			amfLocationRoutes := s.getLocationRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_LOC)
			amfLocationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfLocationGroup, amfLocationRoutes)
		case models.ServiceName_NAMF_OAM:
			amfOAMGroup := router.Group(factory.AmfOamResUriPrefix)
			amfOAMRoutes := s.getOAMRoutes()
			routerAuthorizationCheck := util_oauth.NewRouterAuthorizationCheck(models.ServiceName_NAMF_OAM)
			amfOAMGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, amf_context.GetSelf())
			})
			applyRoutes(amfOAMGroup, amfOAMRoutes)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	var profile models.NfProfile
	if profileTmp, err1 := s.Consumer().BuildNFInstance(s.Context()); err1 != nil {
		logger.InitLog.Error("Build AMF Profile Error")
	} else {
		profile = profileTmp
	}
	_, nfId, err_reg := s.Consumer().SendRegisterNFInstance(s.Context().NrfUri, s.Context().NfId, profile)
	if err_reg != nil {
		logger.InitLog.Warnf("Send Register NF Instance failed: %+v", err_reg)
	} else {
		s.Context().NfId = nfId
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
			s.Terminate()
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
		fmt.Printf("using https")
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetCertPemPath(), cfg.GetCertKeyPath())
	} else {
		err = fmt.Errorf("no support this scheme[%s]", scheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Warnf("SBI server (listen on %s) stopped", s.httpServer.Addr)
}
