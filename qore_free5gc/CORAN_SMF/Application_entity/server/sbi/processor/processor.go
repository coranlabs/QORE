package processor

import (
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_SMF/Application_entity/server/sbi/consumer"
)

const (
	CONTEXT_NOT_FOUND = "CONTEXT_NOT_FOUND"
)

type ProcessorSmf interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorSmf
}

func NewProcessor(smf ProcessorSmf) (*Processor, error) {
	p := &Processor{
		ProcessorSmf: smf,
	}
	return p, nil
}
