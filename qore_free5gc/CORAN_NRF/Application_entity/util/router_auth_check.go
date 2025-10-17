package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_NRF/Application_entity/logger"
	nrf_context "github.com/coranlabs/CORAN_NRF/Messages_handling_entity/context"
)

type (
	NFContextGetter          func() *nrf_context.NRFContext
	RouterAuthorizationCheck struct {
		serviceName models.ServiceName
	}
)

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, nrfContext nrf_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := nrfContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck::Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	logger.UtilLog.Debugf("RouterAuthorizationCheck::Check Authorized")
}
