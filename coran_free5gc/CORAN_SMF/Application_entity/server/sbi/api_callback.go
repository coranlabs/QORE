package sbi

import (
	"net/http"
	"strings"

	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
)

func (s *Server) getCallbackRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/sm-policies/:smContextRef/update",
			APIFunc: s.HTTPSmPolicyUpdateNotification,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/sm-policies/:smContextRef/terminate",
			APIFunc: s.SmPolicyControlTerminationRequestNotification,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/:notifyUri",
			APIFunc: s.HTTPChargingNotification,
		},
	}
}

// SubscriptionsPost -
func (s *Server) HTTPSmPolicyUpdateNotification(c *gin.Context) {
	var request models.SmPolicyNotification

	reqBody, err := c.GetRawData()
	if err != nil {
		logger.PduSessLog.Errorln("GetRawData failed")
	}

	err = openapi.Deserialize(&request, reqBody, c.ContentType())
	if err != nil {
		logger.PduSessLog.Errorln("Deserialize request failed")
	}

	smContextRef := c.Params.ByName("smContextRef")
	s.Processor().HandleSMPolicyUpdateNotify(c, request, smContextRef)
}

func (s *Server) SmPolicyControlTerminationRequestNotification(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPChargingNotification(c *gin.Context) {
	var req models.ChargingNotifyRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.PduSessLog.Errorln("GetRawData failed")
	}

	err = openapi.Deserialize(&req, requestBody, APPLICATION_JSON)
	if err != nil {
		logger.PduSessLog.Errorln("Deserialize request failed")
	}

	smContextRef := strings.Split(c.Params.ByName("notifyUri"), "_")[1]

	s.Processor().HandleChargingNotification(c, req, smContextRef)
}
