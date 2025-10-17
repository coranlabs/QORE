package sbi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coranlabs/CORAN_SCP/Application_entity/factory"
	"github.com/coranlabs/CORAN_SCP/Application_entity/logger"
)

// Structs to match the response structure
type NFServices struct {
	ServiceInstanceId string `json:"serviceInstanceId"`
	ServiceName       string `json:"serviceName"`
	Versions          []struct {
		ApiVersionInUri string `json:"apiVersionInUri"`
		ApiFullVersion  string `json:"apiFullVersion"`
	} `json:"versions"`
	Scheme          string `json:"scheme"`
	NFServiceStatus string `json:"nfServiceStatus"`
	IpEndPoints     []struct {
		Ipv4Address string `json:"ipv4Address"`
		Port        int    `json:"port"`
	} `json:"ipEndPoints"`
}

type NFInstance struct {
	NFInstanceId string `json:"nfInstanceId"`
	NFType       string `json:"nfType"`
	NFStatus     string `json:"nfStatus"`
	PlmnList     []struct {
		Mcc string `json:"mcc"`
		Mnc string `json:"mnc"`
	} `json:"plmnList"`
	Ipv4Addresses []string `json:"ipv4Addresses"`
	AusfInfo      struct{} `json:"ausfInfo"`
	CustomInfo    struct {
		Oauth2 bool `json:"oauth2"`
	} `json:"customInfo"`
	NFServices []NFServices `json:"nfServices"`
}

type NrfResponse struct {
	ValidityPeriod int          `json:"validityPeriod"`
	NFInstances    []NFInstance `json:"nfInstances"`
}

// Function to parse the NRF response and construct the URI
func getURI(responseJSON string, nftype string) (string, error) {
	// Parse the JSON response
	if nftype == "AUSF" {
		var nrfResponse NrfResponse
		err := json.Unmarshal([]byte(responseJSON), &nrfResponse)
		if err != nil {
			return "", fmt.Errorf("error parsing JSON: %v", err)
		}

		// Loop through the NF instances and services to find the AUSF URI
		for _, instance := range nrfResponse.NFInstances {
			if instance.NFType == "AUSF" && instance.NFStatus == "REGISTERED" {
				for _, service := range instance.NFServices {
					if service.ServiceName == "nausf-auth" && service.NFServiceStatus == "REGISTERED" {
						for _, endpoint := range service.IpEndPoints {
							uri := fmt.Sprintf("%s://%s:%d", service.Scheme, endpoint.Ipv4Address, endpoint.Port)
							return uri, nil
						}
					}
				}
			}
		}

		return "", fmt.Errorf("AUSF URI not found in the response")
	} else {
		var result struct {
			NFInstances []struct {
				NFType     string `json:"nfType"`
				NFStatus   string `json:"nfStatus"`
				NFServices []struct {
					ServiceName     string `json:"serviceName"`
					NFServiceStatus string `json:"nfServiceStatus"`
					ApiPrefix       string `json:"apiPrefix"`
				} `json:"nfServices"`
			} `json:"nfInstances"`
		}

		// Parse the JSON response
		err := json.Unmarshal([]byte(responseJSON), &result)
		if err != nil {
			return "", fmt.Errorf("error parsing JSON: %v", err)
		}

		// Look for a registered UDM instance and its services
		for _, instance := range result.NFInstances {
			if instance.NFType == nftype && instance.NFStatus == "REGISTERED" {
				for _, service := range instance.NFServices {
					if service.NFServiceStatus == "REGISTERED" {
						// Return the first found apiPrefix for registered services
						return service.ApiPrefix, nil
					}
				}
			}
		}

		return "", fmt.Errorf("URI not found in the response")
	}
}

func (s *Server) getNFUri(requesternf string, targetnf string) (string, error) {
	var queryPath string
	var url string
	if targetnf == "UDM" {
		queryPath = NrfDiscResUriPrefix + "/nf-instances?requester-nf-type=" + requesternf + "&service-names=nudm-ueau&" + "&target-nf-type=" + targetnf
	} else {
		queryPath = NrfDiscResUriPrefix + "/nf-instances?requester-nf-type=" + requesternf + "&target-nf-type=" + targetnf
	}
	Nrfuri := factory.ScpConfig.GetNrfUri()
	if Nrfuri == "" {
		logger.SBILog.Infof("nrf uri not found in config using default: %s", nrfuri)
		url = nrfuri + queryPath
	} else {
		logger.SBILog.Infof("Nrf uri from config: %s", Nrfuri)
		url = Nrfuri + queryPath
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	// log.Println("respose body : ", string(body))

	uri, err := getURI(string(body), targetnf)
	if err != nil {
		logger.SBILog.Errorf("error: %v", err)
	}

	logger.SBILog.Infof("Requested URI from NRF: %s", uri)
	return uri, nil
	// return string(body), nil
}
