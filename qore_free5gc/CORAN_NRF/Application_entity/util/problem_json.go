package util

import (
	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

func GinProblemJson(c *gin.Context, problemDetails *models.ProblemDetails) {
	c.Writer.Header().Set("Content-Type", "application/problem+json")
	c.JSON(int(problemDetails.Status), problemDetails)
}
