package processor

import (
	"github.com/coranlabs/CORAN_PCF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_PCF/Application_entity/server/sbi/consumer"
)

type PCF interface {
	app.App
	Consumer() *consumer.Consumer
}

type Processor struct {
	PCF
}

func NewProcessor(pcf PCF) (*Processor, error) {
	p := &Processor{
		PCF: pcf,
	}

	return p, nil
}
