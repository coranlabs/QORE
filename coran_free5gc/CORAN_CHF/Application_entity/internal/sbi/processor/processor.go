package processor

import "github.com/coranlabs/CORAN_CHF/Application_entity/pkg/app"

type ProcessorChf interface {
	app.App
}

type Processor struct {
	ProcessorChf
}

type HandlerResponse struct {
	Status  int
	Headers map[string][]string
	Body    interface{}
}

func NewProcessor(chf ProcessorChf) (*Processor, error) {
	p := &Processor{
		ProcessorChf: chf,
	}
	return p, nil
}
