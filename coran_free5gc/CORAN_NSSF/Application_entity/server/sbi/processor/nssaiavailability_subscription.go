/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package processor

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/logger"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_NSSF/Application_entity/util"
)

// Get available subscription ID from configuration
// In this implementation, string converted from 32-bit integer is used as subscription ID
func getUnusedSubscriptionID() (string, error) {
	var idx uint32 = 1
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, subscription := range factory.NssfConfig.Subscriptions {
		tempID, err := strconv.Atoi(subscription.SubscriptionId)
		if err != nil {
			return "", err
		}
		if uint32(tempID) == idx {
			if idx == math.MaxUint32 {
				return "", fmt.Errorf("No available subscription ID")
			}
			idx++
		} else {
			break
		}
	}
	return strconv.Itoa(int(idx)), nil
}

// NSSAIAvailability subscription POST method
func (p *Processor) NssaiAvailabilitySubscriptionCreate(
	c *gin.Context,
	createData models.NssfEventSubscriptionCreateData,
) {
	var (
		response       *models.NssfEventSubscriptionCreatedData = &models.NssfEventSubscriptionCreatedData{}
		problemDetails *models.ProblemDetails
	)

	var subscription factory.Subscription
	tempID, err := getUnusedSubscriptionID()
	if err != nil {
		logger.NssaiavailLog.Warnf(err.Error())

		problemDetails = &models.ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusNotFound,
			Detail: err.Error(),
		}

		util.GinProblemJson(c, problemDetails)
		return
	}

	subscription.SubscriptionId = tempID
	subscription.SubscriptionData = new(models.NssfEventSubscriptionCreateData)
	*subscription.SubscriptionData = createData

	factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions, subscription)

	response.SubscriptionId = subscription.SubscriptionId
	if !subscription.SubscriptionData.Expiry.IsZero() {
		response.Expiry = new(time.Time)
		*response.Expiry = *subscription.SubscriptionData.Expiry
	}
	response.AuthorizedNssaiAvailabilityData = util.AuthorizeOfTaListFromConfig(subscription.SubscriptionData.TaiList)

	c.JSON(http.StatusOK, response)
}

func (p *Processor) NssaiAvailabilitySubscriptionUnsubscribe(c *gin.Context, subscriptionId string) {
	var problemDetails *models.ProblemDetails

	factory.NssfConfig.Lock()
	defer factory.NssfConfig.Unlock()
	for i, subscription := range factory.NssfConfig.Subscriptions {
		if subscription.SubscriptionId == subscriptionId {
			factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions[:i],
				factory.NssfConfig.Subscriptions[i+1:]...)

			c.Status(http.StatusNoContent)
			return
		}
	}

	// No specific subscription ID exists
	problemDetails = &models.ProblemDetails{
		Title:  util.UNSUPPORTED_RESOURCE,
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("Subscription ID '%s' is not available", subscriptionId),
	}

	util.GinProblemJson(c, problemDetails)
}
