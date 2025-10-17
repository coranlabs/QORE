package consumer

import (
	"context"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/Namf_Communication"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFDiscovery"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudr_DataRepository"
	"github.com/coranlabs/CORAN_PCF/Application_entity/pkg/factory"
	pcf_context "github.com/coranlabs/CORAN_PCF/Messages_handling_entity/context"
)

type pcf interface {
	Config() *factory.Config
	Context() *pcf_context.PCFContext
	CancelContext() context.Context
}

type Consumer struct {
	pcf

	// consumer services
	*nnrfService
	*namfService
	*nudrService
}

func NewConsumer(pcf pcf) (*Consumer, error) {
	c := &Consumer{
		pcf: pcf,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
		nfDiscClients:   make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.namfService = &namfService{
		consumer:     c,
		nfComClients: make(map[string]*Namf_Communication.APIClient),
	}

	c.nudrService = &nudrService{
		consumer:         c,
		nfDataSubClients: make(map[string]*Nudr_DataRepository.APIClient),
	}

	return c, nil
}
