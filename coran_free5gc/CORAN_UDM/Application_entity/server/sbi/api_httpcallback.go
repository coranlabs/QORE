package sbi

import (
	"net/http"
	"strings"

	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/coranlabs/CORAN_UDM/Application_entity/logger"
	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

func (s *Server) getHttpCallBackRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			s.HandleIndex,
		},

		{
			"DataChangeNotificationToNF",
			strings.ToUpper("Post"),
			"/sdm-subscriptions",
			s.HandleDataChangeNotificationToNF,
		},
	}
}

func (s *Server) HandleDataChangeNotificationToNF(c *gin.Context) {
	var dataChangeNotify models.DataChangeNotify
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&dataChangeNotify, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	supi := c.Params.ByName("supi")

	logger.CallbackLog.Infof("Handle DataChangeNotificationToNF")

	s.Processor().DataChangeNotificationProcedure(c, dataChangeNotify.NotifyItems, supi)
}
