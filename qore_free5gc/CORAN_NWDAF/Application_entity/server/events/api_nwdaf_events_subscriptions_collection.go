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
 * Description: Routes and config of the events nbi service (create subscription).
 */

package events

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var config EventsConfig

// ------------------------------------------------------------------------------
// Type of EngineConfig structure
type EventsConfig struct {
	Routes struct {
		NumOfUe             string `envconfig:"ENGINE_NUM_OF_UE_ROUTE"`
		SessSuccRatio       string `envconfig:"ENGINE_SESS_SUCC_RATIO_ROUTE"`
		UeComm              string `envconfig:"ENGINE_UE_COMMUNICATION_ROUTE"`
		UeMob               string `envconfig:"ENGINE_UE_MOBILITY_ROUTE"`
		UnexpectedLargeRate string `envconfig:"ENGINE_UNEXPECTED_LARGE_RATE_FLOW_ROUTE"`
	}
	Engine struct {
		Uri    string `envconfig:"ENGINE_URI"`
		AdsUri string `envconfig:"ENGINE_ADS_URI"`
	}
	Events struct {
		Uri string `envconfig:"EVENTS_URI"`
	}
}

// NWDAFEventsSubscriptionsCollectionApiController binds http requests to an api service and writes the service results to the http response
type NWDAFEventsSubscriptionsCollectionApiController struct {
	service      NWDAFEventsSubscriptionsCollectionApiServicer
	errorHandler ErrorHandler
}

// NWDAFEventsSubscriptionsCollectionApiOption for how the controller is set up.
type NWDAFEventsSubscriptionsCollectionApiOption func(
	*NWDAFEventsSubscriptionsCollectionApiController,
)

// ------------------------------------------------------------------------------
// InitConfig - Initialize global variables (config)
func InitConfig() {
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ------------------------------------------------------------------------------
// WithNWDAFEventsSubscriptionsCollectionApiErrorHandler inject ErrorHandler into controller
func WithNWDAFEventsSubscriptionsCollectionApiErrorHandler(
	h ErrorHandler,
) NWDAFEventsSubscriptionsCollectionApiOption {
	return func(c *NWDAFEventsSubscriptionsCollectionApiController) {
		c.errorHandler = h
	}
}

// ------------------------------------------------------------------------------
// NewNWDAFEventsSubscriptionsCollectionApiController creates a default api controller
func NewNWDAFEventsSubscriptionsCollectionApiController(
	s NWDAFEventsSubscriptionsCollectionApiServicer,
	opts ...NWDAFEventsSubscriptionsCollectionApiOption,
) Router {
	controller := &NWDAFEventsSubscriptionsCollectionApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}
	for _, opt := range opts {
		opt(controller)
	}
	return controller
}

// ------------------------------------------------------------------------------
// Routes returns all the api routes for the NWDAFEventsSubscriptionsCollectionApiController
func (c *NWDAFEventsSubscriptionsCollectionApiController) Routes() Routes {
	return Routes{
		{
			"CreateNWDAFEventsSubscription",
			strings.ToUpper("Post"),
			"/nnwdaf-eventssubscription/v1/subscriptions",
			c.CreateNWDAFEventsSubscription,
		},
	}
}

// ------------------------------------------------------------------------------
// CreateNWDAFEventsSubscription - Create a new Individual NWDAF Events Subscription
func (c *NWDAFEventsSubscriptionsCollectionApiController) CreateNWDAFEventsSubscription(
	w http.ResponseWriter,
	r *http.Request,
) {
	nnwdafEventsSubscriptionParam := NnwdafEventsSubscription{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&nnwdafEventsSubscriptionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.CreateNWDAFEventsSubscription(
		r.Context(),
		nnwdafEventsSubscriptionParam,
		r.URL.Path,
	)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)
}
