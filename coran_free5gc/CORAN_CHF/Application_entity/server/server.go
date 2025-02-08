package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coranlabs/CORAN_CHF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_CHF/Application_entity/pkg/service"
)

func Action() error {

	tlsKeyLogPath := ""

	//logger.MainLog.Infoln("CHF version: ", version.GetVersion())

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()

	cfg, err := factory.ReadConfig("./config/CORAN_CHF.yaml")
	if err != nil {
		sigCh <- nil
		time.Sleep(200000)
		return err
	}
	factory.ChfConfig = cfg

	chf, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		sigCh <- nil
		return err
	}
	//CHF = chf

	chf.Start()

	return nil
}
