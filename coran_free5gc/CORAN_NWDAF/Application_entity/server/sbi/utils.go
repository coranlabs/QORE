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
 * Description: This file contains utils functions.
 */

package sbi

import (
	"context"
	"log"
	"time"

	"github.com/coranlabs/CORAN_NWDAF/Application_entity/logger"
	amf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/amfclient"
	smf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/smfclient"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------------------------------------------------------------
// InitConfig - Initialize global variables (cfg and mongoClient) and subscribe to AMF and SMF
func InitConfig() {
	logger.InitializeLogger(logrus.InfoLevel)
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientOptions := options.Client().ApplyURI(config.Database.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.InitLog.Infof("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("error connecting: ", err)
	}
	logger.InitLog.Infof("Connected to MongoDB.")
	mongoClient = client
	// Subscribe to all event notifications from AMF
	amfEventSubscription(
		config.Server.NotifUri+config.Amf.ApiRoute,
		config.Amf.NotifCorrId,
		config.Amf.NotifId,
	)
	// Subscribe to all event notifications from SMF
	smfEventSubscription(
		config.Server.NotifUri+config.Smf.ApiRoute,
		config.Smf.NotifId,
	)
}

// ------------------------------------------------------------------------------
func amfEventSubscription(
	amfEventNotifyUri string,
	amfNotifyCorrelationId string,
	amfNfId string,
) {
	// Store all AMF event types
	var amfEvents []amf_client.AmfEvent
	for _, amfEventTypeAnyOf := range amf_client.AllowedAmfEventTypeAnyOfEnumValues {
		amfEvents = append(amfEvents, *amf_client.NewAmfEvent(amfEventTypeAnyOf))
	}
	// Subscribe to all AMF event types
	amfCreateEventSubscription := *amf_client.NewAmfCreateEventSubscription(
		*amf_client.NewAmfEventSubscription(
			amfEvents,
			amfEventNotifyUri,
			amfNotifyCorrelationId,
			amfNfId,
		),
	)
	configuration := amf_client.NewConfiguration()
	configuration.Debug = true
	amfApiClient := amf_client.NewAPIClient(configuration)
	resp, r, err := amfApiClient.SubscriptionsCollectionCollectionApi.CreateSubscription(
		context.Background()).AmfCreateEventSubscription(amfCreateEventSubscription).Execute()
	if err != nil {
		logger.SubscriberLog.Errorf(
			"Error when calling `SubscriptionsCollectionCollectionApi.CreateSubscription``: %v",
			err,
		)
		log.Printf("Full HTTP response: %v\n", r)
	}
	// response from `CreateSubscription`: AmfCreatedEventSubscription
	logger.SubscriberLog.Infof("Got Response from SubscriptionsCollectionApi AMF")
	logger.SubscriberLog.Infof("Created Individual Subcription for AMF")
	logger.SubscriberLog.Debugf("Response: %v", resp)
}

// ------------------------------------------------------------------------------
func smfEventSubscription(smfEventNotifyUri string, smfNfId string) {

	// Store all SMF event types
	var smfEventSubs []smf_client.EventSubscription
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_PDU_SES_EST),
	)
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_UE_IP_CH),
	)
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_PLMN_CH),
	)
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_DDDS),
	)
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_PDU_SES_REL),
	)
	smfEventSubs = append(smfEventSubs,
		*smf_client.NewEventSubscription(smf_client.SMFEVENTANYOF_QOS_MON),
	)
	// Subscribe to all SMF event types
	nsmfEventExposure := *smf_client.NewNsmfEventExposure(
		smfNfId,
		smfEventNotifyUri,
		smfEventSubs,
	)
	configuration := smf_client.NewConfiguration()
	smfApiClient := smf_client.NewAPIClient(configuration)
	resp, r, err := smfApiClient.SubscriptionsCollectionApi.CreateIndividualSubcription(
		context.Background()).NsmfEventExposure(nsmfEventExposure).Execute()
	if err != nil {
		logger.SubscriberLog.Errorf(
			"Error when calling SubscriptionsCollectionApi.CreateIndividualSubcription: %v",
			err, 
		)
		log.Printf("Full HTTP response: %v\n", r)
	}
	// response from `CreateIndividualSubcription`: NsmfEventExposure
	logger.SubscriberLog.Infof("Got Response from SubscriptionsCollectionApi SMF")
	logger.SubscriberLog.Infof("Created Individual Subcription for SMF")
	logger.SubscriberLog.Debugf("Response: %v", resp)
}
