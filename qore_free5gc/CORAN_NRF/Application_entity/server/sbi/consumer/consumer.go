package consumer

import (
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_NRF/Application_entity/pkg/app"
)

type ConsumerNrf interface {
	app.App
}

type Consumer struct {
	ConsumerNrf

	*nnrfService
}

func NewConsumer(nrf ConsumerNrf) (*Consumer, error) {
	c := &Consumer{
		ConsumerNrf: nrf,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
	}
	return c, nil
}
