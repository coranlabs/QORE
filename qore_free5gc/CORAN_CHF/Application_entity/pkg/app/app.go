package app

import (
	chf_context "github.com/coranlabs/CORAN_CHF/Application_entity/internal/context"
	"github.com/coranlabs/CORAN_CHF/Application_entity/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *chf_context.CHFContext
	Config() *factory.Config
}
