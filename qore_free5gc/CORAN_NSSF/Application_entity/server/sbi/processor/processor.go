package processor

import (
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/app"
)

type Processor struct {
	app.NssfApp
}

func NewProcessor(nssf app.NssfApp) *Processor {
	return &Processor{
		NssfApp: nssf,
	}
}
