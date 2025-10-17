package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	amf_context "github.com/coranlabs/CORAN_AMF/Messages_controller/context"
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

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, amfContext amf_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := amfContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		//logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	//logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Authorized")
}
