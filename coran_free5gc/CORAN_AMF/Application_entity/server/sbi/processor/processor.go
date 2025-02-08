package processor

import (
	"github.com/coranlabs/CORAN_AMF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_AMF/Application_entity/server/sbi/consumer"
)

type ProcessorAmf interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorAmf
}

func NewProcessor(amf ProcessorAmf) (*Processor, error) {
	p := &Processor{
		ProcessorAmf: amf,
	}
	return p, nil
}
