package app

import (
	"github.com/coranlabs/CORAN_NRF/Application_entity/pkg/factory"
	nrf_context "github.com/coranlabs/CORAN_NRF/Messages_handling_entity/context"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *nrf_context.NRFContext
	Config() *factory.Config
}
