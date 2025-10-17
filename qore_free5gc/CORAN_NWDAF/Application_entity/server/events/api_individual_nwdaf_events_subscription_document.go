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
 * Description: Routes and config of the events nbi service (delete, update subscriptions).
 */

package events

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// IndividualNWDAFEventsSubscriptionDocumentApiController binds http requests to an api service and writes the service results to the http response
type IndividualNWDAFEventsSubscriptionDocumentApiController struct {
	service      IndividualNWDAFEventsSubscriptionDocumentApiServicer
	errorHandler ErrorHandler
}

// IndividualNWDAFEventsSubscriptionDocumentApiOption for how the controller is set up.
type IndividualNWDAFEventsSubscriptionDocumentApiOption func(*IndividualNWDAFEventsSubscriptionDocumentApiController)

// WithIndividualNWDAFEventsSubscriptionDocumentApiErrorHandler inject ErrorHandler into controller
func WithIndividualNWDAFEventsSubscriptionDocumentApiErrorHandler(h ErrorHandler) IndividualNWDAFEventsSubscriptionDocumentApiOption {
	return func(c *IndividualNWDAFEventsSubscriptionDocumentApiController) {
		c.errorHandler = h
	}
}

// NewIndividualNWDAFEventsSubscriptionDocumentApiController creates a default api controller
func NewIndividualNWDAFEventsSubscriptionDocumentApiController(s IndividualNWDAFEventsSubscriptionDocumentApiServicer, opts ...IndividualNWDAFEventsSubscriptionDocumentApiOption) Router {
	controller := &IndividualNWDAFEventsSubscriptionDocumentApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the IndividualNWDAFEventsSubscriptionDocumentApiController
func (c *IndividualNWDAFEventsSubscriptionDocumentApiController) Routes() Routes {
	return Routes{
		{
			"DeleteNWDAFEventsSubscription",
			strings.ToUpper("Delete"),
			"/nnwdaf-eventssubscription/v1/subscriptions/{subscriptionId}",
			c.DeleteNWDAFEventsSubscription,
		},
		{
			"UpdateNWDAFEventsSubscription",
			strings.ToUpper("Put"),
			"/nnwdaf-eventssubscription/v1/subscriptions/{subscriptionId}",
			c.UpdateNWDAFEventsSubscription,
		},
	}
}

// DeleteNWDAFEventsSubscription - Delete an existing Individual NWDAF Events Subscription
func (c *IndividualNWDAFEventsSubscriptionDocumentApiController) DeleteNWDAFEventsSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	subscriptionIdParam := params["subscriptionId"]
	result, err := c.service.DeleteNWDAFEventsSubscription(r.Context(), subscriptionIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)
}

// UpdateNWDAFEventsSubscription - Update an existing Individual NWDAF Events Subscription
func (c *IndividualNWDAFEventsSubscriptionDocumentApiController) UpdateNWDAFEventsSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	subscriptionIdParam := params["subscriptionId"]
	nnwdafEventsSubscriptionParam := NnwdafEventsSubscription{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&nnwdafEventsSubscriptionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertNnwdafEventsSubscriptionRequired(nnwdafEventsSubscriptionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateNWDAFEventsSubscription(r.Context(), subscriptionIdParam, nnwdafEventsSubscriptionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, result.Headers, w)
}
