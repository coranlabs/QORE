package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// "github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/service"
	// "github.com/coranlabs/CORAN_LIB_UTIL/version"
)

var UDR *service.UdrApp

func Action() error {
	tlsKeyLogPath := ""

	// logger.MainLog.Infoln("UDR version: ", version.GetVersion())
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	cfg, err := factory.ReadConfig("./config/CORAN_UDR.yaml")
	if err != nil {
		return err
	}
	factory.UdrConfig = cfg
	udr, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}
	UDR = udr

	udr.Start()

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
// 		_, name := filepath.Split(factory.UdrDefaultTLSKeyLogPath)
// 		logTlsKeyPath = filepath.Join(tmpDir, name)
// 	}

// 	return logTlsKeyPath, nil
// }
