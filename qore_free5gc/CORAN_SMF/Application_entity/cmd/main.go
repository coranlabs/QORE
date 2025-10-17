package cmd

import (
	"context"
	// "math/rand"
	"os"
	"os/signal"
	"path/filepath"

	// "runtime/debug"
	"syscall"
	// "time"

	// "github.com/urfave/cli"

	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	// "github.com/coranlabs/CORAN_LIB_UTIL/version"
	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/service"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/utils"
)

var SMF *service.SmfApp

func action() error {
	tlsKeyLogPath := ""

	// logger.MainLog.Infoln("SMF version: ", version.GetVersion())

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	cfg, err := factory.ReadConfig("./config/CORAN_SMF.yaml")
	if err != nil {
		sigCh <- nil
		return err
	}
	factory.SmfConfig = cfg

	ueRoutingCfg, err := factory.ReadUERoutingConfig("./config/UEROUTING.yaml")
	if err != nil {
		sigCh <- nil
		return err
	}
	factory.UERoutingConfig = ueRoutingCfg

	pfcpStart, pfcpTerminate := utils.InitPFCPFunc()
	smf, err := service.NewApp(ctx, cfg, tlsKeyLogPath, pfcpStart, pfcpTerminate)
	if err != nil {
		sigCh <- nil
		return err
	}
	SMF = smf

	smf.Start()

	return nil
}

func initLogFile(logNfPath []string) (string, error) {
	logTlsKeyPath := ""

	for _, path := range logNfPath {
		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
			return "", err
		}

		if logTlsKeyPath != "" {
			continue
		}

		nfDir, _ := filepath.Split(path)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return "", err
		}
		_, name := filepath.Split(factory.SmfDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}

	return logTlsKeyPath, nil
}
