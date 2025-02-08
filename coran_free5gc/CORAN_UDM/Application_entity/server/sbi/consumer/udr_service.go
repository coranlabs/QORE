package consumer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nudr_DataRepository"
	"github.com/coranlabs/CORAN_UDM/Application_entity/logger"
	udm_context "github.com/coranlabs/CORAN_UDM/Message_controller/context"
)

type nudrService struct {
	consumer *Consumer

	nfDRMu sync.RWMutex

	nfDRClients map[string]*Nudr_DataRepository.APIClient
}

const (
	NFDiscoveryToUDRParamNone int = iota
	NFDiscoveryToUDRParamSupi
	NFDiscoveryToUDRParamExtGroupId
	NFDiscoveryToUDRParamGpsi
)

func (s *nudrService) CreateUDMClientToUDR(id string) (*Nudr_DataRepository.APIClient, error) {
	uri := s.getUdrURI(id)
	if uri == "" {
		logger.ProcLog.Errorf("ID[%s] does not match any UDR", id)
		return nil, fmt.Errorf("No UDR URI found")
	}
	s.nfDRMu.RLock()
	client, ok := s.nfDRClients[uri]
	if ok {
		s.nfDRMu.RUnlock()
		return client, nil
	}

	cfg := Nudr_DataRepository.NewConfiguration()
	cfg.SetBasePath(uri)
	client = Nudr_DataRepository.NewAPIClient(cfg)

	s.nfDRMu.RUnlock()
	s.nfDRMu.Lock()
	defer s.nfDRMu.Unlock()
	s.nfDRClients[uri] = client
	return client, nil
}

func (s *nudrService) getUdrURI(id string) string {
	if strings.Contains(id, "imsi") || strings.Contains(id, "nai") { // supi
		ue, ok := udm_context.GetSelf().UdmUeFindBySupi(id)
		if ok {
			if ue.UdrUri == "" {
				ue.UdrUri = SendNFIntancesUDR(id, NFDiscoveryToUDRParamSupi)
			}
			return ue.UdrUri
		} else {
			ue = udm_context.GetSelf().NewUdmUe(id)
			ue.UdrUri = SendNFIntancesUDR(id, NFDiscoveryToUDRParamSupi)
			return ue.UdrUri
		}
	} else if strings.Contains(id, "pei") {
		var udrURI string
		udm_context.GetSelf().UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.Amf3GppAccessRegistration != nil && ue.Amf3GppAccessRegistration.Pei == id {
				if ue.UdrUri == "" {
					ue.UdrUri = SendNFIntancesUDR(ue.Supi, NFDiscoveryToUDRParamSupi)
				}
				udrURI = ue.UdrUri
				return false
			} else if ue.AmfNon3GppAccessRegistration != nil && ue.AmfNon3GppAccessRegistration.Pei == id {
				if ue.UdrUri == "" {
					ue.UdrUri = SendNFIntancesUDR(ue.Supi, NFDiscoveryToUDRParamSupi)
				}
				udrURI = ue.UdrUri
				return false
			}
			return true
		})
		return udrURI
	} else if strings.Contains(id, "extgroupid") {
		// extra group id
		return SendNFIntancesUDR(id, NFDiscoveryToUDRParamExtGroupId)
	} else if strings.Contains(id, "msisdn") || strings.Contains(id, "extid") {
		// gpsi
		return SendNFIntancesUDR(id, NFDiscoveryToUDRParamGpsi)
	}
	return SendNFIntancesUDR("", NFDiscoveryToUDRParamNone)
}
