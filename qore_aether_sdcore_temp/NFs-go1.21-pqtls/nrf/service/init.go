// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/lakshya-chopra/http2_util"
	mongoDBLibLogger "github.com/omec-project/MongoDBLibrary/logger"
	"github.com/omec-project/logger_util"
	"github.com/omec-project/nrf/accesstoken"
	nrf_context "github.com/omec-project/nrf/context"
	"github.com/omec-project/nrf/dbadapter"
	"github.com/omec-project/nrf/discovery"
	"github.com/omec-project/nrf/factory"
	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/nrf/management"
	"github.com/omec-project/nrf/util"
	"github.com/omec-project/path_util"
	pathUtilLogger "github.com/omec-project/path_util/logger"
)

type NRF struct{}

type (
	// Config information.
	Config struct {
		nrfcfg string
	}
)

var config Config

var nrfCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "free5gccfg",
		Usage: "common config file",
	},
	cli.StringFlag{
		Name:  "nrfcfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NRF) GetCliCmd() (flags []cli.Flag) {
	return nrfCLi
}

func (nrf *NRF) Initialize(c *cli.Context) error {
	config = Config{
		nrfcfg: c.String("nrfcfg"),
	}

	if config.nrfcfg != "" {
		if err := factory.InitConfigFactory(config.nrfcfg); err != nil {
			return err
		}
	} else {
		DefaultNrfConfigPath := path_util.Free5gcPath("free5gc/config/nrfcfg.conf")
		if err := factory.InitConfigFactory(DefaultNrfConfigPath); err != nil {
			return err
		}
	}

	nrf.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nrf *NRF) setLogLevel() {
	if factory.NrfConfig.Logger == nil {
		initLog.Warnln("NRF config without log level setting!!!")
		return
	}

	if factory.NrfConfig.Logger.NRF != nil {
		if factory.NrfConfig.Logger.NRF.DebugLevel != "" {
			level, err := logrus.ParseLevel(factory.NrfConfig.Logger.NRF.DebugLevel)
			if err != nil {
				initLog.Warnf("NRF Log level [%s] is invalid, set to [info] level",
					factory.NrfConfig.Logger.NRF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("NRF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("NRF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NrfConfig.Logger.NRF.ReportCaller)
	}

	if factory.NrfConfig.Logger.PathUtil != nil {
		if factory.NrfConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NrfConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.NrfConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.NrfConfig.Logger.PathUtil.ReportCaller)
	}

	/*if factory.NrfConfig.Logger.OpenApi != nil {
		if factory.NrfConfig.Logger.OpenApi.DebugLevel != "" {
			if _, err := logrus.ParseLevel(factory.NrfConfig.Logger.OpenApi.DebugLevel); err != nil {
				logger.OpenapiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.NrfConfig.Logger.OpenApi.DebugLevel)
			}
		} else {
			logger.OpenapiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
		}
		logger.SetReportCaller(factory.NrfConfig.Logger.OpenApi.ReportCaller)
	}*/

	if factory.NrfConfig.Logger.MongoDBLibrary != nil {
		if factory.NrfConfig.Logger.MongoDBLibrary.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NrfConfig.Logger.MongoDBLibrary.DebugLevel); err != nil {
				mongoDBLibLogger.MongoDBLog.Warnf("MongoDBLibrary Log level [%s] is invalid, set to [info] level",
					factory.NrfConfig.Logger.MongoDBLibrary.DebugLevel)
				mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				mongoDBLibLogger.SetLogLevel(level)
			}
		} else {
			mongoDBLibLogger.MongoDBLog.Warnln("MongoDBLibrary Log level not set. Default set to [info] level")
			mongoDBLibLogger.SetLogLevel(logrus.InfoLevel)
		}
		mongoDBLibLogger.SetReportCaller(factory.NrfConfig.Logger.MongoDBLibrary.ReportCaller)
	}
}

func (nrf *NRF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nrf.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
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

	fmt.Printf("%s End %s\n", sep, sep)
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

func (nrf *NRF) Start() {
	initLog.Infoln("Server started")
	dbadapter.ConnectToDBClient(factory.NrfConfig.Configuration.MongoDBName, factory.NrfConfig.Configuration.MongoDBUrl,
		factory.NrfConfig.Configuration.MongoDBStreamEnable, factory.NrfConfig.Configuration.NfProfileExpiryEnable)

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	accesstoken.AddService(router)
	discovery.AddService(router)
	management.AddService(router)

	nrf_context.InitNrfContext()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		// Waiting for other NFs to deregister
		time.Sleep(2 * time.Second)
		nrf.Terminate()
		os.Exit(0)
	}()

	roc := os.Getenv("MANAGED_BY_CONFIG_POD")
	if roc == "true" {
		initLog.Infoln("MANAGED_BY_CONFIG_POD is true")
	} else {
		initLog.Infoln("Use helm chart config ")
	}
	bindAddr := factory.NrfConfig.GetSbiBindingAddr()
	initLog.Infof("Binding addr: [%s]", bindAddr)

	cert, err := tls.LoadX509KeyPair(util.NrfPemPath, util.NrfKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	//print
	cert_x509, _ := ReadCertificate(util.NrfPemPath)
	PrintCertificateDetails(cert_x509)

	server, err := http2_util.NewServer(bindAddr, util.NrfLogPath, router, cert)

	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: +%v", err)
	}

	serverScheme := factory.NrfConfig.GetSbiScheme()
	fmt.Printf("\nServer scheme: %s\n", serverScheme)
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS("", "")
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (nrf *NRF) Exec(c *cli.Context) error {
	initLog.Traceln("args:", c.String("nrfcfg"))
	args := nrf.FilterCli(c)
	initLog.Traceln("filter: ", args)
	command := exec.Command("./nrf", args...)

	if err := nrf.Initialize(c); err != nil {
		return err
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		initLog.Fatalln(err)
	}
	go func() {
		in := bufio.NewScanner(stderr)
		fmt.Println("NRF log start")
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		fmt.Println("NRF  start")
		if err = command.Start(); err != nil {
			fmt.Printf("NRF Start error: %v", err)
		}
		fmt.Println("NRF  end")
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (nrf *NRF) Terminate() {
	logger.InitLog.Infof("Terminating NRF...")

	logger.InitLog.Infof("NRF terminated")
}
