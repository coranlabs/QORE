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
 * Description: Functions of the events nbi service (delete, update subscriptions).
 */

package events

import (
	"context"
	"errors"
	"net/http"
)

// IndividualNWDAFEventsSubscriptionDocumentApiService is a service that implements the logic for the IndividualNWDAFEventsSubscriptionDocumentApiServicer
// This service should implement the business logic for every endpoint for the IndividualNWDAFEventsSubscriptionDocumentApi API.
// Include any external packages or services that will be required by this service.
type IndividualNWDAFEventsSubscriptionDocumentApiService struct {
}

// NewIndividualNWDAFEventsSubscriptionDocumentApiService creates a default api service
func NewIndividualNWDAFEventsSubscriptionDocumentApiService() IndividualNWDAFEventsSubscriptionDocumentApiServicer {
	return &IndividualNWDAFEventsSubscriptionDocumentApiService{}
}

// DeleteNWDAFEventsSubscription - Delete an existing Individual NWDAF Events Subscription
func (s *IndividualNWDAFEventsSubscriptionDocumentApiService) DeleteNWDAFEventsSubscription(
	ctx context.Context,
	subscriptionId string,
) (ImplResponse, error) {
	subscriptionCh, ok := subscriptionTable[subscriptionId]
	// If the key exists
	if ok {
		close(subscriptionCh)
		delete(subscriptionTable, subscriptionId)

		return Response(204, nil), nil
	}
	return Response(http.StatusBadRequest, ProblemDetails{}), nil
}

// UpdateNWDAFEventsSubscription - Update an existing Individual NWDAF Events Subscription
func (s *IndividualNWDAFEventsSubscriptionDocumentApiService) UpdateNWDAFEventsSubscription(
	ctx context.Context,
	subscriptionId string,
	nnwdafEventsSubscription NnwdafEventsSubscription,
) (ImplResponse, error) {
	// TODO - update UpdateNWDAFEventsSubscription with the required logic for this service method.
	// Add api_individual_nwdaf_events_subscription_document_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, NnwdafEventsSubscription{}) or use other options such as http.Ok ...
	//return Response(200, NnwdafEventsSubscription{}), nil

	//TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	//return Response(204, nil),nil

	//TODO: Uncomment the next line to return response Response(307, RedirectResponse{}) or use other options such as http.Ok ...
	//return Response(307, RedirectResponse{}), nil

	//TODO: Uncomment the next line to return response Response(308, RedirectResponse{}) or use other options such as http.Ok ...
	//return Response(308, RedirectResponse{}), nil

	//TODO: Uncomment the next line to return response Response(400, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(400, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(401, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(401, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(403, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(403, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(404, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(404, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(411, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(411, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(413, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(413, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(415, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(415, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(429, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(429, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(500, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(500, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(501, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(501, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(503, ProblemDetails{}) or use other options such as http.Ok ...
	//return Response(503, ProblemDetails{}), nil

	//TODO: Uncomment the next line to return response Response(0, {}) or use other options such as http.Ok ...
	//return Response(0, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("UpdateNWDAFEventsSubscription method not implemented")
}
