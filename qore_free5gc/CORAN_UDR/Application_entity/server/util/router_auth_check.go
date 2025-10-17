package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	udr_context "github.com/coranlabs/CORAN_UDR/Message_controller/context"
	// "github.com/coranlabs/CORAN_UDR/Application_entity/logger"
)

type RouterAuthorizationCheck struct {
	serviceName models.ServiceName
}

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, udrContext udr_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := udrContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Authorized")
}
