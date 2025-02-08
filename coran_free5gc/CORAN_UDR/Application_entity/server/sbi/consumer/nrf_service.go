package consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFDiscovery"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/Nnrf_NFManagement"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	udr_context "github.com/coranlabs/CORAN_UDR/Message_controller/context"
)

type NrfService struct {
	nfMngmntMu sync.RWMutex

	nfMngmntClients map[string]*Nnrf_NFManagement.APIClient
}

func (ns *NrfService) getNFManagementClient(uri string) *Nnrf_NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	ns.nfMngmntMu.RLock()
	client, ok := ns.nfMngmntClients[uri]
	if ok {
		ns.nfMngmntMu.RUnlock()
		return client
	}

	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnrf_NFManagement.NewAPIClient(configuration)

	ns.nfMngmntMu.RUnlock()
	ns.nfMngmntMu.Lock()
	defer ns.nfMngmntMu.Unlock()
	ns.nfMngmntClients[uri] = client
	return client
}

func (ns *NrfService) buildNFProfile(context *udr_context.UDRContext) (models.NfProfile, error) {
	// config := factory.UdrConfig

	profile := models.NfProfile{
		NfInstanceId:  context.NfId,
		NfType:        models.NfType_UDR,
		NfStatus:      models.NfStatus_REGISTERED,
		Ipv4Addresses: []string{context.RegisterIPv4},
		UdrInfo: &models.UdrInfo{
			SupportedDataSets: []models.DataSetId{
				// models.DataSetId_APPLICATION,
				// models.DataSetId_EXPOSURE,
				// models.DataSetId_POLICY,
				models.DataSetId_SUBSCRIPTION,
			},
		},
	}

	var services []models.NfService
	for _, nfService := range context.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = &services
	}

	return profile, nil
}

func (ns *NrfService) SendRegisterNFInstance(ctx context.Context, nrfUri string) (string, string, error) {
	// Set client and set url
	profile, err := ns.buildNFProfile(udr_context.GetSelf())
	if err != nil {
		return "", "", fmt.Errorf("failed to build nrf profile %s", err.Error())
	}

	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := ns.getNFManagementClient(nrfUri)
	var resouceNrfUri string
	var retrieveNfInstanceId string

	finish := false

	for !finish {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context done")
		default:
			nf, res, registerErr := client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, profile.NfInstanceId, profile)
			if registerErr != nil || res == nil {
				// TODO : add log
				logger.ConsumerLog.Errorf("UDR register to NRF Error[%s]", registerErr.Error())
				time.Sleep(2 * time.Second)
				continue
			}
			defer func() {
				if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
					logger.ConsumerLog.Errorf("RegisterNFInstance response body cannot close: %+v", rspCloseErr)
				}
			}()

			status := res.StatusCode
			if status == http.StatusOK {
				// NFUpdate
				finish = true
			} else if status == http.StatusCreated {
				// NFRegister
				resourceUri := res.Header.Get("Location")
				resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
				retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

				oauth2 := false
				if nf.CustomInfo != nil {
					v, ok := nf.CustomInfo["oauth2"].(bool)
					if ok {
						oauth2 = v
						logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
					}
				}
				udr_context.GetSelf().OAuth2Required = oauth2
				if oauth2 && udr_context.GetSelf().NrfCertPem == "" {
					logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
				}
				finish = true
			} else {
				logger.ConsumerLog.Errorf("NRF returned wrong status code: %d", status)
			}
		}
	}
	return resouceNrfUri, retrieveNfInstanceId, nil
}

func (ns *NrfService) SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	ctx, pd, err := udr_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NfType_NRF)
	if err != nil {
		return pd, err
	}

	udrSelf := udr_context.GetSelf()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(udrSelf.NrfUri)
	client := ns.getNFManagementClient(udrSelf.NrfUri)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, udrSelf.NfId)
	if err == nil {
		return nil, err
	} else if res != nil {
		defer func() {
			if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
				logger.ConsumerLog.Errorf("DeregisterNFInstance response body cannot close: %+v", rspCloseErr)
			}
		}()

		if res.Status != err.Error() {
			return nil, err
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}
	return problemDetails, err
}

func (ns *NrfService) SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts,
) (*models.SearchResult, error) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	ctx, _, err := udr_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NfType_NRF)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, targetNfType, requestNfType, &param)
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("temporary redirect for non NRF consumer")
	}
	defer func() {
		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
			logger.ConsumerLog.Errorf("SearchNFInstances response body cannot close: %+v", rspCloseErr)
		}
	}()

	return &result, err
}
