/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the OAI Public License, Version 1.1  (the "License"); you may not use this
 * file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 *      http://www.openairinterface.org/?page_id=698
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*
 * Author: Abdelkader Mekrache <mekrache@eurecom.fr>
 * Author: Karim Boutiba 	   <boutiba@eurecom.fr>
 * Author: Arina Prostakova    <prostako@eurecom.fr>
 * Description: Routes and config of the analytics nbi service.
 */

package analytics

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var config AnalyticsConfig

// ------------------------------------------------------------------------------
// Type of EngineConfig structure
type AnalyticsConfig struct {
	Routes struct {
		NumOfUe       string `envconfig:"ENGINE_NUM_OF_UE_ROUTE"`
		SessSuccRatio string `envconfig:"ENGINE_SESS_SUCC_RATIO_ROUTE"`
		UeComm        string `envconfig:"ENGINE_UE_COMMUNICATION_ROUTE"`
		UeMob         string `envconfig:"ENGINE_UE_MOBILITY_ROUTE"`
	}
	Engine struct {
		Uri string `envconfig:"ENGINE_URI"`
	}
}

// ------------------------------------------------------------------------------
type NWDAFAnalyticsDocumentApiController struct {
	service      NWDAFAnalyticsDocumentApiServicer
	errorHandler ErrorHandler
}

// NWDAFAnalyticsDocumentApiOption for how the controller is set up.
type NWDAFAnalyticsDocumentApiOption func(*NWDAFAnalyticsDocumentApiController)

// ------------------------------------------------------------------------------
// InitConfig - Initialize global variables (config)
func InitConfig() {
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ------------------------------------------------------------------------------
// WithNWDAFAnalyticsDocumentApiErrorHandler inject ErrorHandler into controller
func WithNWDAFAnalyticsDocumentApiErrorHandler(
	h ErrorHandler,
) NWDAFAnalyticsDocumentApiOption {
	return func(c *NWDAFAnalyticsDocumentApiController) {
		c.errorHandler = h
	}
}

// ------------------------------------------------------------------------------
// NewNWDAFAnalyticsDocumentApiController creates a default api controller
func NewNWDAFAnalyticsDocumentApiController(
	s NWDAFAnalyticsDocumentApiServicer,
	opts ...NWDAFAnalyticsDocumentApiOption,
) Router {
	controller := &NWDAFAnalyticsDocumentApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}
	for _, opt := range opts {
		opt(controller)
	}
	return controller
}

// ------------------------------------------------------------------------------
// Routes returns all the api routes for the NWDAFAnalyticsDocumentApiController
func (c *NWDAFAnalyticsDocumentApiController) Routes() Routes {
	return Routes{
		{
			"GetNWDAFAnalytics",
			strings.ToUpper("Get"),
			"/nnwdaf-analyticsinfo/v1/analytics",
			c.GetNWDAFAnalytics,
		},
	}
}

// ------------------------------------------------------------------------------
// GetNWDAFAnalytics - Read a NWDAF Analytics
func (c *NWDAFAnalyticsDocumentApiController) GetNWDAFAnalytics(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("Getting NWDAF Analytics")
	query := r.URL.Query()
	eventIdParam := query.Get("event-id")
	anaReqParam := query.Get("ana-req")
	var anaReq EventReportingRequirement
	json.Unmarshal([]byte(anaReqParam), &anaReq)
	eventFilterParam := query.Get("event-filter")
	var eventFilter EventFilter
	json.Unmarshal([]byte(eventFilterParam), &eventFilter)
	supportedFeaturesParam := query.Get("supported-features")
	tgtUeParam := query.Get("tgt-ue")
	var tgtUe TargetUeInformation
	json.Unmarshal([]byte(tgtUeParam), &tgtUe)
	result, err := c.service.GetNWDAFAnalytics(
		r.Context(),
		EventIdAnyOf(eventIdParam),
		anaReq,
		eventFilter,
		supportedFeaturesParam,
		tgtUe,
	)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
