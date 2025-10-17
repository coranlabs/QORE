package sbi

import (
	"log"

	// "github.com/coranlabs/CORAN_LIB_OPENAPI/models"

	"github.com/gin-gonic/gin"
)

const (
	AmfCallbackResUriPrefix = "/namf-callback/v1"
	AmfCommResUriPrefix     = "/namf-comm/v1"
	AmfEvtsResUriPrefix     = "/namf-evts/v1"
	AmfLocResUriPrefix      = "/namf-loc/v1"
	AmfMtResUriPrefix       = "/namf-mt/v1"
	AmfOamResUriPrefix      = "/namf-oam/v1"

	AusfSorprotectionResUriPrefix = "/nausf-sorprotection/v1"
	AusfAuthResUriPrefix          = "/nausf-auth/v1"
	AusfUpuprotectionResUriPrefix = "/nausf-upuprotection/v1"

	NssfNssaiavailResUriPrefix = "/nnssf-nssaiavailability/v1"
	NssfNsselectResUriPrefix   = "/nnssf-nsselection/v1"

	PcfPolicyAuthResUriPrefix   = "/npcf-policyauthorization/v1"
	PcfAMpolicyCtlResUriPrefix  = "/npcf-am-policy-control/v1"
	PcfCallbackResUriPrefix     = "/npcf-callback/v1"
	PcfSMpolicyCtlResUriPrefix  = "/npcf-smpolicycontrol/v1"
	PcfBdtPolicyCtlResUriPrefix = "/npcf-bdtpolicycontrol/v1"
	PcfOamResUriPrefix          = "/npcf-oam/v1"
	PcfUePolicyCtlResUriPrefix  = "/npcf-ue-policy-control/v1/"

	SmfEventExposureResUriPrefix = "/nsmf_event-exposure/v1"
	SmfPdusessionResUriPrefix    = "/nsmf-pdusession/v1"
	SmfOamUriPrefix              = "/nsmf-oam/v1"
	SmfCallbackUriPrefix         = "/nsmf-callback"
	// UpiUriPrefix                 = "/upi/v1"

	UdmSorprotectionResUriPrefix  = "/nudm-sorprotection/v1"
	UdmAuthResUriPrefix           = "/nudm-auth/v1"
	UdmfUpuprotectionResUriPrefix = "/nudm-upuprotection/v1"
	UdmEcmResUriPrefix            = "/nudm-ecm/v1"
	UdmSdmResUriPrefix            = "/nudm-sdm/v1"
	UdmEeResUriPrefix             = "/nudm-ee/v1"
	// UdmDrResUriPrefix             = "/nudr-dr/v1"
	UdmUecmResUriPrefix = "/nudm-uecm/v1"
	UdmPpResUriPrefix   = "/nudm-pp/v1"
	UdmUeauResUriPrefix = "/nudm-ueau/v1"

	UdrDrResUriPrefix = "/nudr-dr/v1"

	NrfNfmResUriPrefix  = "/nnrf-nfm/v1"
	NrfDiscResUriPrefix = "/nnrf-disc/v1"

	ConvergedChargingResUriPrefix = "/nchf-convergedcharging/v3"
)

const (
	amfuri   = "http://amf.free5gc.org:8000"
	ausfuri  = "http://coran-free5gc-ausf-service:80" //"http://ausf.free5gc.org:8000"
	pcfuri   = "http://pcf.free5gc.org:8000"
	smfuri   = "http://smf.free5gc.org:8000"
	udmuri   = "http://udm.free5gc.org:8000"
	udruri   = "http://udr.free5gc.org:8000"
	nssfuri  = "http://nssf.free5gc.org:8000"
	n3iwfuri = "http://n3iwf.free5gc.org:8000"
	nrfuri   = "http://nrf-nnrf:8000"
)

type NFUriResult struct {
	Uri string
	Err error
}

// func (s *Server) getNFUri2() string {

// 	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}
// 	resp, err := consumer.GetConsumer().SendSearchNFInstances(
// 		NrfUri, models.NfType_AUSF, models.NfType_AMF, &param)
// 	if err != nil {
// 		// ue.GmmLog.Error("AMF can not select an AUSF by NRF")
// 		// gmm_message.SendRegistrationReject(ue.RanUe[accessType], nasMessage.Cause5GMMCongestion, "")
// 		return ""
// 	}
// 	var ausfUri string
// 	for _, nfProfile := range resp.NfInstances {
// 		// AusfId := nfProfile.NfInstanceId
// 		ausfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NAUSF_AUTH, models.NfServiceStatus_REGISTERED)
// 		if ausfUri != "" {
// 			break
// 		}
// 	}
// 	return ausfUri
// }

// Handler function that checks the group and forwards requests accordingly
func (s *Server) handleNausfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	// Ausfuri, err := s.getNFUri("SCP", "AUSF")
	resultChan := make(chan NFUriResult)

	// Run the getNFUri function in a Goroutine
	go func() {
		uri, err := s.getNFUri("SCP", "AUSF")
		if err != nil {
			log.Println(err)
		}
		// Send the result to the channel
		resultChan <- NFUriResult{Uri: uri, Err: err}
	}()

	result := <-resultChan
	// log.Println("ausf uri from nrf: ", result.Uri)
	Ausfuri := result.Uri
	// log.Println("ausf uri from nrf: ", Ausfuri)
	s.forwardRequestAusf(c, Ausfuri)
}

func (s *Server) handleNpcfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	Pcfuri, err := s.getNFUri("SCP", "PCF")
	if err != nil {
		log.Fatalln(err)
	}
	s.forwardRequestPcf(c, Pcfuri)
}

func (s *Server) handleNudmRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	// Udmuri, err := s.getNFUri("SCP", "UDM")
	resultChan := make(chan NFUriResult)

	// Run the getNFUri function in a Goroutine
	go func() {
		uri, err := s.getNFUri("SCP", "UDM")
		if err != nil {
			log.Println(err)
		}
		// Send the result to the channel
		resultChan <- NFUriResult{Uri: uri, Err: err}
	}()

	result := <-resultChan
	// log.Println("ausf uri from nrf: ", result.Uri)
	Udmuri := result.Uri
	s.forwardRequestUdm(c, Udmuri)
}

func (s *Server) handleNamfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	Amfuri, err := s.getNFUri("SCP", "AMF")
	if err != nil {
		log.Fatalln(err)
	}
	s.forwardRequestAmf(c, Amfuri)
}

func (s *Server) handleNsmfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	Smfuri, err := s.getNFUri("SCP", "SMF")
	if err != nil {
		log.Fatalln(err)
	}
	s.forwardRequestSmf(c, Smfuri)
}
func (s *Server) handleNnssfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	// Nssfuri, err := s.getNFUri("SCP", "NSSF")
	resultChan := make(chan NFUriResult)

	// Run the getNFUri function in a Goroutine
	go func() {
		uri, err := s.getNFUri("SCP", "NSSF")
		if err != nil {
			log.Println(err)
		}
		// Send the result to the channel
		resultChan <- NFUriResult{Uri: uri, Err: err}
	}()

	result := <-resultChan
	// log.Println("uri from nrf: ", result.Uri)
	Nssfuri := result.Uri
	s.forwardRequestNssf(c, Nssfuri)
}
func (s *Server) handleNudrRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	// Udruri, err := s.getNFUri("SCP", "UDR")
	resultChan := make(chan NFUriResult)

	// Run the getNFUri function in a Goroutine
	go func() {
		uri, err := s.getNFUri("SCP", "UDR")
		if err != nil {
			log.Println(err)
		}
		// Send the result to the channel
		resultChan <- NFUriResult{Uri: uri, Err: err}
	}()

	result := <-resultChan
	Udruri := result.Uri
	s.forwardRequestUdr(c, Udruri)
}

func (s *Server) handleChfRoutes(c *gin.Context) {
	// service := c.Param("service") // Get the service type from the URL
	// action := c.Param("action")   // Get the action part of the path
	// log.Printf("Service: %s, Action: %s", service, action)
	// Udruri, err := s.getNFUri("SCP", "UDR")
	resultChan := make(chan NFUriResult)

	// Run the getNFUri function in a Goroutine
	go func() {
		uri, err := s.getNFUri("SCP", "CHF")
		if err != nil {
			log.Println(err)
		}
		// Send the result to the channel
		resultChan <- NFUriResult{Uri: uri, Err: err}
	}()

	result := <-resultChan
	Udruri := result.Uri
	s.forwardRequestChf(c, Udruri)
}
func (s *Server) handleNn3iwfRoutes(c *gin.Context) {
	service := c.Param("service") // Get the service type from the URL
	action := c.Param("action")   // Get the action part of the path
	log.Printf("Service: %s, Action: %s", service, action)
	// forwardRequest(c, n3iwfuri)
}
