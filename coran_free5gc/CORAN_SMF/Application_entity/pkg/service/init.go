package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/consumer"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/processor"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
)

type SmfAppInterface interface {
	app.App

	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

var SMF SmfAppInterface

type SmfApp struct {
	SmfAppInterface

	cfg    *factory.Config
	smfCtx *smf_context.SMFContext
	ctx    context.Context
	cancel context.CancelFunc

	sbiServer *sbi.Server
	consumer  *consumer.Consumer
	processor *processor.Processor
	wg        sync.WaitGroup

	pfcpStart     func(*SmfApp)
	pfcpTerminate func()
}

func GetApp() SmfAppInterface {
	return SMF
}

func NewApp(
	ctx context.Context, cfg *factory.Config, tlsKeyLogPath string,
	pfcpStart func(*SmfApp), pfcpTerminate func(),
) (*SmfApp, error) {
	smf_context.Init()
	smf := &SmfApp{
		cfg:           cfg,
		wg:            sync.WaitGroup{},
		pfcpStart:     pfcpStart,
		pfcpTerminate: pfcpTerminate,
		smfCtx:        smf_context.GetSelf(),
	}
	smf.SetLogEnable(cfg.GetLogEnable())
	smf.SetLogLevel(cfg.GetLogLevel())
	smf.SetReportCaller(cfg.GetLogReportCaller())

	// Initialize consumer
	consumer, err := consumer.NewConsumer(smf)
	if err != nil {
		return nil, err
	}
	smf.consumer = consumer

	// Initialize processor
	processor, err := processor.NewProcessor(smf)
	if err != nil {
		return nil, err
	}
	smf.processor = processor

	// TODO: Initialize sbi server
	sbiServer, err := sbi.NewServer(smf, tlsKeyLogPath)
	if err != nil {
		return nil, err
	}
	smf.sbiServer = sbiServer

	smf.ctx, smf.cancel = context.WithCancel(ctx)

	// for PFCP
	ctx, cancel := context.WithCancel(smf.ctx)
	smf_context.GetSelf().Ctx = ctx
	smf_context.GetSelf().PFCPCancelFunc = cancel

	SMF = smf

	return smf, nil
}

func (a *SmfApp) Config() *factory.Config {
	return a.cfg
}

func (a *SmfApp) Context() *smf_context.SMFContext {
	return a.smfCtx
}

func (a *SmfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *SmfApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *SmfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *SmfApp) SetLogEnable(enable bool) {
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

func (a *SmfApp) SetLogLevel(level string) {
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

func (a *SmfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *SmfApp) Start() {
	logger.InitLog.Infoln("Server started")

	err := a.sbiServer.Run(context.Background(), &a.wg)
	if err != nil {
		logger.MainLog.Errorf("sbi server run error %+v", err)
	}

	a.wg.Add(1)
	go a.listenShutDownEvent()

	// Initialize PFCP server
	a.pfcpStart(a)

	a.WaitRoutineStopped()
}

func (a *SmfApp) listenShutDownEvent() {
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

func (a *SmfApp) Terminate() {
	a.cancel()
}

func (a *SmfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating SMF...")
	a.pfcpTerminate()
	// deregister with NRF
	problemDetails, err := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.MainLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.MainLog.Errorf("Deregister NF instance Error[%+v]", err)
	} else {
		logger.MainLog.Infof("Deregister from NRF successfully")
	}
	a.sbiServer.Stop()
	logger.MainLog.Infof("SMF SBI Server terminated")
}

func (a *SmfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("SMF App is terminated")
}
