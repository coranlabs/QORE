package app

import (
	"github.com/coranlabs/CORAN_AUSF/Application_entity/pkg/factory"
	ausf_context "github.com/coranlabs/CORAN_AUSF/Messages_handling_entity/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *ausf_context.AUSFContext
	Config() *factory.Config
}
