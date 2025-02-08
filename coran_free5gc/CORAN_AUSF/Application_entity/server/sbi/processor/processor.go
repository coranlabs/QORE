package processor

import (
	"github.com/coranlabs/CORAN_AUSF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_AUSF/Application_entity/server/sbi/consumer"
)

type ProcessorAusf interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorAusf
}

func NewProcessor(ausf ProcessorAusf) (*Processor, error) {
	p := &Processor{
		ProcessorAusf: ausf,
	}
	return p, nil
}
