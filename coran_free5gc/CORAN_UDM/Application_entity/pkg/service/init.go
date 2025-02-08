package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/coranlabs/CORAN_UDM/Application_entity/config/factory"
	"github.com/coranlabs/CORAN_UDM/Application_entity/logger"
	"github.com/coranlabs/CORAN_UDM/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_UDM/Application_entity/server/sbi"
	"github.com/coranlabs/CORAN_UDM/Application_entity/server/sbi/consumer"
	"github.com/coranlabs/CORAN_UDM/Application_entity/server/sbi/processor"
	"github.com/sirupsen/logrus"

	udm_context "github.com/coranlabs/CORAN_UDM/Message_controller/context"
)

var _ app.App = &UdmApp{}

type UdmApp struct {
	udmCtx *udm_context.UDMContext
	cfg    *factory.Config

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	sbiServer *sbi.Server
	consumer  *consumer.Consumer
	processor *processor.Processor
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*UdmApp, error) {
	udm := &UdmApp{
		cfg: cfg,
		wg:  sync.WaitGroup{},
	}
	udm.SetLogEnable(cfg.GetLogEnable())
	udm.SetLogLevel(cfg.GetLogLevel())
	udm.SetReportCaller(cfg.GetLogReportCaller())
	udm_context.Init()

	consumer, err := consumer.NewConsumer(udm)
	if err != nil {
		return udm, err
	}
	udm.consumer = consumer

	processor, err_p := processor.NewProcessor(udm)
	if err_p != nil {
		return udm, err_p
	}
	udm.processor = processor

	udm.ctx, udm.cancel = context.WithCancel(ctx)
	udm.udmCtx = udm_context.GetSelf()

	if udm.sbiServer, err = sbi.NewServer(udm, tlsKeyLogPath); err != nil {
		return nil, err
	}

	return udm, nil
}

func (a *UdmApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *UdmApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *UdmApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}
	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *UdmApp) Start() {
	logger.InitLog.Infoln("Server started")

	a.wg.Add(1)
	go a.listenShutdownEvent()

	if err := a.sbiServer.Run(context.Background(), &a.wg); err != nil {
		logger.MainLog.Fatalf("Run SBI server failed: %+v", err)
	}

	a.WaitRoutineStopped()
}

func (a *UdmApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.terminateProcedure()
}

func (a *UdmApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Stop()
	}
}

func (a *UdmApp) Terminate() {
	a.cancel()
}

func (a *UdmApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating UDM...")
	a.CallServerStop()

	// deregister with NRF
	problemDetails, err := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.MainLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.MainLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.MainLog.Infof("Deregister from NRF successfully")
	}
	logger.MainLog.Infof("UDM SBI Server terminated")
}

func (a *UdmApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("UDM App is terminated")
}

func (a *UdmApp) Config() *factory.Config {
	return a.cfg
}

func (a *UdmApp) Context() *udm_context.UDMContext {
	return a.udmCtx
}

func (a *UdmApp) CancelContext() context.Context {
	return a.ctx
}

func (a *UdmApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *UdmApp) Processor() *processor.Processor {
	return a.processor
}
