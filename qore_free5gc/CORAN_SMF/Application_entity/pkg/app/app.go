package app

import (
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *smf_context.SMFContext
	Config() *factory.Config
}
