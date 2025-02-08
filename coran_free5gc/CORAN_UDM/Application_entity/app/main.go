package cmd

import (
	"context"
	// "fmt"
	// "net"
	"os"
	"os/signal"
	"syscall"

	"github.com/coranlabs/CORAN_UDM/Application_entity/config/factory"
	"github.com/coranlabs/CORAN_UDM/Application_entity/logger"
	"github.com/coranlabs/CORAN_UDM/Application_entity/pkg/service"
	// "github.com/CORAN_UDM/Application_entity/server"
	// "github.com/gin-gonic/gin"
)

// ListenAndServe starts a TCP listener on the specified address
func Action() error {
	tlsKeyLogPath := ""
	// if err != nil {
	// 	return err
	// }

	// logger.MainLog.Infoln("UDM version: ", version.GetVersion())
	logger.AppLog.Infof("UDM started")
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	cfg, err := factory.ReadConfig("./config/CORAN_UDM.yaml")
	if err != nil {
		return err
	}
	factory.UdmConfig = cfg

	udm, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}

	udm.Start()

	return nil
}
