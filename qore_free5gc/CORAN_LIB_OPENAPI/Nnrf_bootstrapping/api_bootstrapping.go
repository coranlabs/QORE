package Nnrf_Bootstrapping

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	// "github.com/antihax/optional"

	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

type BootstrappingRequestApiService service

func (b *BootstrappingRequestApiService) BootstrappingRequest(ctx context.Context) (models.BootstrappingInfo, *http.Response, error) {
	var (
		localVarHTTPMethod   = strings.ToUpper("Get")
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  models.BootstrappingInfo
	)

	// create path and map variables
	localVarPath := b.client.cfg.BasePath() + "/bootstrapping"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	localVarHTTPContentTypes := []string{"application/json"}
	localVarHeaderParams["Content-Type"] = localVarHTTPContentTypes[0]

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}
	localVarHTTPHeaderAccept := openapi.SelectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}

	r, err := openapi.PrepareRequest(ctx, b.client.cfg, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := openapi.CallAPI(b.client.cfg, r)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	apiError := openapi.GenericOpenAPIError{
		RawBody:     localVarBody,
		ErrorStatus: localVarHTTPResponse.Status,
	}

	switch localVarHTTPResponse.StatusCode {
	case 200:
		err = openapi.Deserialize(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			apiError.ErrorStatus = err.Error()
		}
		return localVarReturnValue, localVarHTTPResponse, nil
	default:
		return localVarReturnValue, localVarHTTPResponse, openapi.ReportError("%d is not a valid status code in BootstrappingRequest", localVarHTTPResponse.StatusCode)
	}
}
