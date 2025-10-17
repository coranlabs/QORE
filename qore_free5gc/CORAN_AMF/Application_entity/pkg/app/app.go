package app

import (
	"github.com/coranlabs/CORAN_AMF/Application_entity/config/factory"
	amf_context "github.com/coranlabs/CORAN_AMF/Messages_controller/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *amf_context.AMFContext
	Config() *factory.Config
}
