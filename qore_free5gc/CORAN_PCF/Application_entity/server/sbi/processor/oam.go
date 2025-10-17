package processor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_PCF/Application_entity/logger"
	"github.com/coranlabs/CORAN_PCF/Messages_handling_entity/context"
)

type UEAmPolicy struct {
	PolicyAssociationID string
	AccessType          models.AccessType
	Rfsp                string
	Triggers            []models.RequestTrigger
	/*Service Area Restriction */
	RestrictionType models.RestrictionType
	Areas           []models.Area
	MaxNumOfTAs     int32
}

type UEAmPolicys []UEAmPolicy

func (p *Processor) HandleOAMGetAmPolicyRequest(
	c *gin.Context,
	supi string,
) {
	// step 1: log
	logger.OamLog.Infof("Handle OAMGetAmPolicy")

	// step 3: handle the message

	logger.OamLog.Infof("Handle OAM Get Am Policy")
	response := &UEAmPolicys{}
	pcfSelf := context.GetSelf()

	if val, exists := pcfSelf.UePool.Load(supi); exists {
		ue := val.(*context.UeContext)
		for _, amPolicy := range ue.AMPolicyData {
			ueAmPolicy := UEAmPolicy{
				PolicyAssociationID: amPolicy.PolAssoId,
				AccessType:          amPolicy.AccessType,
				Rfsp:                strconv.Itoa(int(amPolicy.Rfsp)),
				Triggers:            amPolicy.Triggers,
			}
			if amPolicy.ServAreaRes != nil {
				servAreaRes := amPolicy.ServAreaRes
				ueAmPolicy.RestrictionType = servAreaRes.RestrictionType
				ueAmPolicy.Areas = servAreaRes.Areas
				ueAmPolicy.MaxNumOfTAs = servAreaRes.MaxNumOfTAs
			}
			*response = append(*response, ueAmPolicy)
		}
		c.JSON(http.StatusOK, response)
		return
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
}
