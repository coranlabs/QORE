package consumer

import (
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Namf_Communication"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nchf_ConvergedCharging"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFDiscovery"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Npcf_SMPolicyControl"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nsmf_PDUSession"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_SubscriberDataManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudm_UEContextManagement"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/app"
)

type Consumer struct {
	app.App

	// consumer services
	*nsmfService
	*namfService
	*nchfService
	*npcfService
	*nudmService
	*nnrfService
}

func NewConsumer(smf app.App) (*Consumer, error) {
	c := &Consumer{
		App: smf,
	}

	c.nsmfService = &nsmfService{
		consumer:          c,
		PDUSessionClients: make(map[string]*Nsmf_PDUSession.APIClient),
	}

	c.namfService = &namfService{
		consumer:             c,
		CommunicationClients: make(map[string]*Namf_Communication.APIClient),
	}

	c.nchfService = &nchfService{
		consumer:                 c,
		ConvergedChargingClients: make(map[string]*Nchf_ConvergedCharging.APIClient),
	}

	c.nudmService = &nudmService{
		consumer:                        c,
		SubscriberDataManagementClients: make(map[string]*Nudm_SubscriberDataManagement.APIClient),
		UEContextManagementClients:      make(map[string]*Nudm_UEContextManagement.APIClient),
	}

	c.nnrfService = &nnrfService{
		consumer:            c,
		NFManagementClients: make(map[string]*Nnrf_NFManagement.APIClient),
		NFDiscoveryClients:  make(map[string]*Nnrf_NFDiscovery.APIClient),
	}

	c.npcfService = &npcfService{
		consumer:               c,
		SMPolicyControlClients: make(map[string]*Npcf_SMPolicyControl.APIClient),
	}

	return c, nil
}
