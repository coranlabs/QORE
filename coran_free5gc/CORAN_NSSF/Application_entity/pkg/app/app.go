package app

import (
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/factory"
	nssf_context "github.com/coranlabs/CORAN_NSSF/Messages_handling_entity/context"
)

type NssfApp interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *nssf_context.NSSFContext
	Config() *factory.Config
}
