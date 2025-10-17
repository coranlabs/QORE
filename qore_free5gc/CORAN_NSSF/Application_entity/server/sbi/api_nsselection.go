package sbi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_NSSF/Application_entity/logger"
)

func (s *Server) getNsSelectionRoutes() []Route {
	return []Route{
		{
			"Health Check",
			strings.ToUpper("Get"),
			"/",
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"status": "Service Available"})
			},
		},

		{
			"NSSelectionGet",
			strings.ToUpper("Get"),
			"/network-slice-information",
			s.NetworkSliceInformationGet,
		},
	}
}

func (s *Server) NetworkSliceInformationGet(c *gin.Context) {
	logger.NsselLog.Infof("Handle NSSelectionGet")

	query := c.Request.URL.Query()
	s.Processor().NSSelectionSliceInformationGet(c, query)
}
