package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_AUSF/Application_entity/logger"
	ausf_context "github.com/coranlabs/CORAN_AUSF/Messages_handling_entity/context"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

type RouterAuthorizationCheck struct {
	serviceName models.ServiceName
}

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, ausfContext ausf_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := ausfContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Authorized")
}
