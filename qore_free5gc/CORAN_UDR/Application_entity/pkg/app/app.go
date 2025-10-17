package app

import (
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
	udr_context "github.com/coranlabs/CORAN_UDR/Message_controller/context"
	// "github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *udr_context.UDRContext
	Config() *factory.Config
}
