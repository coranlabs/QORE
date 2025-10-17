package sbi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/coranlabs/CORAN_NRF/Application_entity/logger"
	"github.com/coranlabs/CORAN_NRF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_NRF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_NRF/Application_entity/server/sbi/processor"
	"github.com/coranlabs/CORAN_NRF/Application_entity/util"
	"github.com/lakshya-chopra/httpwrapper"
)

type ServerNrf interface {
	app.App

	// Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

type Server struct {
	ServerNrf

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

func NewServer(nrf ServerNrf, tlsKeyLogPath string) (*Server, error) {
	s := &Server{
		ServerNrf: nrf,
		router:    logger_util.NewGinWithLogrus(logger.GinLog),
	}
	cfg := s.Config()
	bindAddr := cfg.GetSbiBindingAddr()
	logger.SBILog.Infof("Binding addr: [%s]", bindAddr)

	s.applyService()

	server_cert, err1 := tls.LoadX509KeyPair(cfg.GetNrfCertPemPath(), cfg.GetNrfPrivKeyPath())

	if err1 != nil {
		log.Fatal(err1)
	}

	cert, err2 := ReadCertificate(cfg.GetNrfCertPemPath())
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
	return s, nil
}

func (s *Server) GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.NfmLog.Error(err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (s *Server) applyService() {
	accesstokenRoutes := s.getAccesstokenRoutes()
	accesstokenGroup := s.router.Group("") // accesstoken service didn't have api prefix
	applyRoutes(accesstokenGroup, accesstokenRoutes)

	discoveryRoutes := s.getNfDiscoveryRoutes()
	discoveryGroup := s.router.Group(factory.NrfDiscResUriPrefix)
	discAuthCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NNRF_DISC)
	discoveryGroup.Use(func(c *gin.Context) {
		discAuthCheck.Check(c, s.Context())
	})
	applyRoutes(discoveryGroup, discoveryRoutes)

	// OAuth2 must exclude NfRegister
	nfRegisterRoute := s.getNfRegisterRoute()
	nfRegisterGroup := s.router.Group(factory.NrfNfmResUriPrefix)
	applyRoutes(nfRegisterGroup, nfRegisterRoute)

	managementRoutes := s.getNfManagementRoute()
	managementGroup := s.router.Group(factory.NrfNfmResUriPrefix)
	managementAuthCheck := util.NewRouterAuthorizationCheck(models.ServiceName_NNRF_NFM)
	managementGroup.Use(func(c *gin.Context) {
		managementAuthCheck.Check(c, s.Context())
	})
	applyRoutes(managementGroup, managementRoutes)
}

func (s *Server) Run(wg *sync.WaitGroup) error {
	wg.Add(1)
	go s.startServer(wg)

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

	cfg := s.Config()
	serverScheme := cfg.GetSbiScheme()

	var err error
	if serverScheme == "http" {
		err = s.httpServer.ListenAndServe()
	} else if serverScheme == "https" {
		// TODO: support TLS mutual authentication for OAuth
		err = s.httpServer.ListenAndServeTLS(
			cfg.GetNrfCertPemPath(),
			cfg.GetNrfPrivKeyPath())
	} else {
		err = fmt.Errorf("No support this scheme[%s]", serverScheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.SBILog.Errorf("SBI server error: %v", err)
	}
	logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
}

func (s *Server) Stop() {
	// server stop
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(toCtx); err != nil {
		logger.SBILog.Errorf("Could not close SBI server: %#v", err)
	}
}
