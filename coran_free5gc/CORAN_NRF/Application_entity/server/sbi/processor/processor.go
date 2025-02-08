package processor

import (
	"github.com/coranlabs/CORAN_NRF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_NRF/Application_entity/server/sbi/consumer"
)

type ProcessorNrf interface {
	app.App
	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorNrf
}

func NewProcessor(nrf ProcessorNrf) (*Processor, error) {
	p := &Processor{
		ProcessorNrf: nrf,
	}
	return p, nil
}
