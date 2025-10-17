package app

import (
	"github.com/coranlabs/CORAN_PCF/Application_entity/pkg/factory"
	pcf_context "github.com/coranlabs/CORAN_PCF/Messages_handling_entity/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *pcf_context.PCFContext
	Config() *factory.Config
}
