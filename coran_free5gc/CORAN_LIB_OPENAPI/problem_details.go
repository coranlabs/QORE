package CORAN_LIB_OPENAPI

import (
	"net/http"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

func ProblemDetailsSystemFailure(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "System failure",
		Status: http.StatusInternalServerError,
		Detail: detail,
		Cause:  "SYSTEM_FAILURE",
	}
}

func ProblemDetailsMalformedReqSyntax(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "Malformed request syntax",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func ProblemDetailsDataNotFound(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "Data not found",
		Status: http.StatusNotFound,
		Detail: detail,
	}
}
