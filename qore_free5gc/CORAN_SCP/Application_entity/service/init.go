package service

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/coranlabs/CORAN_SCP/Application_entity/factory"
	"github.com/coranlabs/CORAN_SCP/Application_entity/logger"
	"github.com/coranlabs/CORAN_SCP/Application_entity/server/sbi"
	"github.com/sirupsen/logrus"
)

var SCP *ScpApp

type ScpApp struct {
	cfg *factory.Config

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	sbiServer *sbi.Server
}

// Config implements sbi.ServerScp.
func (s *ScpApp) Config() *factory.Config {
	return s.cfg
}

// SetLogEnable implements sbi.ServerScp.
func (s *ScpApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	s.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

// SetLogLevel implements sbi.ServerScp.
func (s *ScpApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	s.cfg.SetLogLevel(level)
	logger.InitializeLogger(lvl)
}

// SetReportCaller implements sbi.ServerScp.
func (s *ScpApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	s.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

// Start implements sbi.ServerScp.

func (s *ScpApp) Start() {
	logger.InitLog.Infoln("Server starting")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	s.wg.Add(1)
	go s.listenShutdownEvent()

	if err := s.sbiServer.Run(&s.wg); err != nil {
		logger.MainLog.Fatalf("Run SBI server failed: %+v", err)
	}

	for {
		select {
		case <-sigChan:
			logger.InitLog.Infoln("Received shutdown signal")
			s.wg.Wait()
			logger.InitLog.Infoln("Server stopped")
			return
		}
	}
}

func (s *ScpApp) listenShutdownEvent() {
	// Create a channel to listen for OS signals (like SIGINT/SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if p := recover(); p != nil {
			// Handle panic and print stack trace to logs
			// logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		s.wg.Done() // Mark the goroutine as done
	}()

	// Block until a signal is received or context is canceled
	select {
	case <-sigChan: // Ctrl+C or SIGTERM signal
		logger.MainLog.Infoln("Received termination signal (Ctrl+C or SIGTERM)")
	case <-s.ctx.Done(): // Context is canceled
		logger.MainLog.Infoln("Context canceled")
	}

	// Initiate the termination procedure
	s.terminateProcedure()
}

func (s *ScpApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating SCP...")

	waitTime := 5
	logger.MainLog.Infof("Waiting for %vs for other NFs to deregister", waitTime)
	// s.waitNfDeregister(waitTime)

	logger.MainLog.Infof("Remove NF Profile...")
	s.sbiServer.Stop()
}

// Terminate implements sbi.ServerScp.
func (s *ScpApp) Terminate() {
	panic("unimplemented")
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*ScpApp, error) {
	scp := &ScpApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}

	scp.SetLogEnable(cfg.GetLogEnable())
	scp.SetLogLevel(cfg.GetLogLevel())
	scp.SetReportCaller(cfg.GetLogReportCaller())
	err := dummy()
	if scp.sbiServer, err = sbi.NewServer(scp, tlsKeyLogPath); err != nil {
		return nil, err
	}
	SCP = scp

	return scp, nil
}

func dummy() error {
	return nil
}
