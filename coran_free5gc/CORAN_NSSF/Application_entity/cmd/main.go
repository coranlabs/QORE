package cmd

import (
	"context"
	"os"
	"os/signal"

	// "path/filepath"

	"syscall"

	// logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
	"github.com/coranlabs/CORAN_LIB_UTIL/version"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/logger"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/service"
)

var NSSF *service.NssfApp

func Action() error {
	tlsKeyLogPath := ""
	logger.MainLog.Infoln("NSSF version: ", version.GetVersion())
	cfg, err := factory.ReadConfig("./config/CORAN_NSSF.yaml")
	if err != nil {
		return err
	}
	factory.NssfConfig = cfg

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	nssf, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}
	NSSF = nssf

	nssf.Start()

	return nil
}

// func initLogFile(logNfPath []string) (string, error) {
// 	logTlsKeyPath := ""

// 	for _, path := range logNfPath {
// 		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
// 			return "", err
// 		}

// 		if logTlsKeyPath != "" {
// 			continue
// 		}
// 		nfDir, _ := filepath.Split(path)
// 		tmpDir := filepath.Join(nfDir, "key")
// 		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
// 			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
// 			return "", err
// 		}
// 		_, name := filepath.Split(factory.NssfDefaultTLSKeyLogPath)
// 		logTlsKeyPath = filepath.Join(tmpDir, name)
// 	}

// 	return logTlsKeyPath, nil
// }
