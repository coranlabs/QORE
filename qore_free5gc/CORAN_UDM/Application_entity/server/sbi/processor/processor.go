package processor

import (
	"github.com/coranlabs/CORAN_UDM/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_UDM/Application_entity/server/sbi/consumer"
)

type ProcessorUdm interface {
	app.App

	Consumer() *consumer.Consumer
}

type Processor struct {
	ProcessorUdm
}

func NewProcessor(udm ProcessorUdm) (*Processor, error) {
	p := &Processor{
		ProcessorUdm: udm,
	}
	return p, nil
}
