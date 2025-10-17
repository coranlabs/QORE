package consumer

import (
	"github.com/coranlabs/CORAN_AMF/Application_entity/pkg/app"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Namf_Communication"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nausf_UEAuthentication"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFDiscovery"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnssf_NSSelection"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Npcf_AMPolicy"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nsmf_PDUSession"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_SubscriberDataManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_UEContextManagement"
)

var consumer *Consumer

type ConsumerAmf interface {
	app.App
}

type Consumer struct {
	ConsumerAmf

	// consumer services
	*namfService
	*nnrfService
	*npcfService
	*nssfService
	*nsmfService
	*nudmService
	*nausfService
}

func GetConsumer() *Consumer {
	return consumer
}

func NewConsumer(amf ConsumerAmf) (*Consumer, error) {
	c := &Consumer{
		ConsumerAmf: amf,
	}

	c.namfService = &namfService{
		consumer:   c,
		ComClients: make(map[string]*Namf_Communication.APIClient),
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfMngmntClients: make(map[string]*Nnrf_NFManagement.APIClient),
		nfDiscClients:   make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.npcfService = &npcfService{
		consumer:        c,
		AMPolicyClients: make(map[string]*Npcf_AMPolicy.APIClient),
	}

	c.nssfService = &nssfService{
		consumer:           c,
		NSSelectionClients: make(map[string]*Nnssf_NSSelection.APIClient),
	}

	c.nsmfService = &nsmfService{
		consumer:          c,
		PDUSessionClients: make(map[string]*Nsmf_PDUSession.APIClient),
	}

	c.nudmService = &nudmService{
		consumer:                 c,
		SubscriberDMngmntClients: make(map[string]*Nudm_SubscriberDataManagement.APIClient),
		UEContextMngmntClients:   make(map[string]*Nudm_UEContextManagement.APIClient),
	}

	c.nausfService = &nausfService{
		consumer:                c,
		UEAuthenticationClients: make(map[string]*Nausf_UEAuthentication.APIClient),
	}
	consumer = c
	return c, nil
}
