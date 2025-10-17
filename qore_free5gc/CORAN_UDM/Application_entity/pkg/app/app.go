package app

import (
	"github.com/coranlabs/CORAN_UDM/Application_entity/config/factory"

	udm_context "github.com/coranlabs/CORAN_UDM/Message_controller/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *udm_context.UDMContext
	Config() *factory.Config
}
