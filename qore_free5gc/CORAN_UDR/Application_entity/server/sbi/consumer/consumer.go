package consumer

import (
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/app"
)

type Consumer struct {
	app.App

	*NrfService
}

func NewConsumer(udr app.App) *Consumer {
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(udr.Context().NrfUri)
	nrfService := &NrfService{
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
	}

	return &Consumer{
		App:        udr,
		NrfService: nrfService,
	}
}
