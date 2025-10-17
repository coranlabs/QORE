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
 * Description: This file contains data structures and global variables.
 */

package sbi

import (
	amf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/amfclient"
	smf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/smfclient"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global variable
var mongoClient *mongo.Client
var config SbiConfig

// ------------------------------------------------------------------------------
// Type of EngineConfig structure
type SbiConfig struct {
	Amf struct {
		IpAddr            string `envconfig:"AMF_IP_ADDR"`
		SubRoute          string `envconfig:"AMF_SUBSCR_ROUTE"`
		ApiRoute          string `envconfig:"AMF_API_ROUTE"`
		NotifCorrId       string `envconfig:"AMF_NOTIFY_CORRELATION_ID"`
		NotifId           string `envconfig:"AMF_NOTIFICATION_ID"`
		NorifForwardRoute string `envconfig:"AMF_NOTIFICATION_FORWARD_ROUTE"`
	}
	Smf struct {
		IpAddr            string `envconfig:"SMF_IP_ADDR"`
		SubRoute          string `envconfig:"SMF_SUBSCR_ROUTE"`
		ApiRoute          string `envconfig:"SMF_API_ROUTE"`
		NotifCorrId       string `envconfig:"SMF_NOTIFY_CORRELATION_ID"`
		NotifId           string `envconfig:"SMF_NOTIFICATION_ID"`
		NorifForwardRoute string `envconfig:"SMF_NOTIFICATION_FORWARD_ROUTE"`
	}
	Database struct {
		Uri               string `envconfig:"MONGODB_URI"`
		DbName            string `envconfig:"MONGODB_DATABASE_NAME"`
		CollectionAmfName string `envconfig:"MONGODB_COLLECTION_NAME_AMF"`
		CollectionSmfName string `envconfig:"MONGODB_COLLECTION_NAME_SMF"`
	}
	Server struct {
		NotifUri string `envconfig:"EVENT_NOTIFY_URI"`
		Uri      string `envconfig:"SERVER_ADDR"`
	}
}
type pduSesEst struct {
	AdIpv4Addr  *string
	Dnn         *string
	PduSeId     *int32
	PduSessType *smf_client.PduSessionType
	Snssai      *smf_client.Snssai
	TimeStamp   int64
}

type ueIpCh struct {
	AdIpv4Addr *string
	PduSeId    *int32
	TimeStamp  int64
}

type ddds struct {
	DddStatus *smf_client.DlDataDeliveryStatus
	PduSeId   *int32
	TimeStamp int64
}

type qosMon struct {
	Customized_data *smf_client.CustomizedData
	PduSeId         *int32
	TimeStamp       int64
}

type rmInfo struct {
	RmInfo    amf_client.RmInfo
	TimeStamp int64
}

type location struct {
	UserLocation amf_client.UserLocation
	TimeStamp    int64
}

type lossOfConnectReason struct {
	LossOfConnectReason amf_client.LossOfConnectivityReasonAnyOf
	TimeStamp           int64
}
