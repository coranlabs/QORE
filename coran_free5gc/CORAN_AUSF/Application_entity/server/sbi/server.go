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

	"github.com/coranlabs/CORAN_AUSF/Application_entity/logger"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/server/sbi/consumer"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/server/sbi/processor"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/util"
	ausf_context "github.com/coranlabs/CORAN_AUSF/Messages_handling_entity/context"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/lakshya-chopra/httpwrapper"
)

type ServerAusf interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerAusf

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

func NewServer(ausf ServerAusf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerAusf: ausf,
	}

	s.router = newRouter(s)

	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)

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

	for _, serverName := range factory.AusfConfig.Configuration.ServiceNameList {
		switch models.ServiceName(serverName) {
		case models.ServiceName_NAUSF_AUTH:
			ausfUeAuthenticationGroup := router.Group(factory.AusfAuthResUriPrefix)
			ausfUeAuthenticationRoutes := s.getUeAuthenticationRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_AUTH)
			ausfUeAuthenticationGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfUeAuthenticationGroup, ausfUeAuthenticationRoutes)
		case models.ServiceName_NAUSF_SORPROTECTION:
			ausfSorprotectionGroup := router.Group(factory.AusfSorprotectionResUriPrefix)
			ausfSorprotectionRoutes := s.getSorprotectionRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_SORPROTECTION)
			ausfSorprotectionGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfSorprotectionGroup, ausfSorprotectionRoutes)
		case models.ServiceName_NAUSF_UPUPROTECTION:
			ausfUpuprotectionGroup := router.Group(factory.AusfUpuprotectionResUriPrefix)
			ausfUpuprotectionRoutes := s.getUpuprotectionRoutes()
			routerAuthorizationCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NAUSF_UPUPROTECTION)
			ausfUpuprotectionGroup.Use(func(c *gin.Context) {
				routerAuthorizationCheck.Check(c, ausf_context.GetSelf())
			})
			applyRoutes(ausfUpuprotectionGroup, ausfUpuprotectionRoutes)
		}
	}

	return router
}

func (s *Server) Run(traceCtx context.Context, wg *sync.WaitGroup) error {
	var err error
	_, s.Context().NfId, err = s.Consumer().RegisterNFInstance(context.Background())
	if err != nil {
		logger.InitLog.Errorf("AUSF register to NRF Error[%s]", err.Error())
	}

	wg.Add(1)
	go s.startServer(wg)

	return nil
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
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
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetCertPemPath(),
			cfg.GetCertKeyPath())
	} else {
		err = fmt.Errorf("No support this scheme[%s]", scheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}
