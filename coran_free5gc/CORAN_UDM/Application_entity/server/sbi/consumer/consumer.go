package consumer

import (
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFDiscovery"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_SubscriberDataManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_UEContextManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudr_DataRepository"
	"github.com/coranlabs/CORAN_UDM/Application_entity/pkg/app"
)

type ConsumerUdm interface {
	app.App
}

type Consumer struct {
	ConsumerUdm

	// consumer services
	*nnrfService
	*nudrService
	*nudmService
}

func NewConsumer(udm ConsumerUdm) (*Consumer, error) {
	c := &Consumer{
		ConsumerUdm: udm,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
		nfDiscClients:   make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.nudrService = &nudrService{
		consumer:    c,
		nfDRClients: make(map[string]*Nudr_DataRepository.APIClient),
	}

	c.nudmService = &nudmService{
		consumer:      c,
		nfSDMClients:  make(map[string]*Nudm_SubscriberDataManagement.APIClient),
		nfUECMClients: make(map[string]*Nudm_UEContextManagement.APIClient),
	}
	return c, nil
}
