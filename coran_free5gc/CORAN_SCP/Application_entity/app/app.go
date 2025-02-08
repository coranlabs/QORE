package app

import "github.com/coranlabs/CORAN_SCP/Application_entity/factory"

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Config() *factory.Config
}
