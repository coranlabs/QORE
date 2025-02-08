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
 * Description: This file contains functions related smf post notifications.
 */

package sbi

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	smf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/smfclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------------------------------------------------------------
func storeSmfNotificationOnDB(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST":
		log.Printf("Storing SMF notification in Database")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var smfNotification *smf_client.NsmfEventExposureNotification
		err = json.Unmarshal(body, &smfNotification)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}
		notifList, ok := smfNotification.GetEventNotifsOk()
		if !ok {
			http.Error(w, "Error in getting EventNotifs from SMF notification", http.StatusBadRequest)
			return
		}
		databaseName := config.Database.DbName
		collectionName := config.Database.CollectionSmfName
		smfCollection := mongoClient.Database(databaseName).Collection(collectionName)
		opts := options.Update().SetUpsert(true)
		// store reports one by one
		for _, notif := range notifList {
			oid := notif.GetSupi()
			if oid == "" {
				http.Error(w, "supi not found in notification, cannot create object id", http.StatusBadRequest)
				return
			}
			update, err := getUpdateByNotif(notif)
			if err != nil {
				http.Error(w, "error in getUpdateByNotif", http.StatusBadRequest)
				return
			}
			// Update/Insert the SMF notification report
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			res, err := smfCollection.UpdateByID(ctx, oid, update, opts)
			if err != nil {
				http.Error(w, "error in updating the SMF collection", http.StatusBadRequest)
				return
			}
			if res.MatchedCount != 0 {
				log.Printf("matched and updated an existing notification report from SMF")
			}
			if res.UpsertedCount != 0 {
				log.Printf("inserted a new notification report from SMF with ID %v\n",
					res.UpsertedID)
			}
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------
// getUpdateByNotif - Return update bson.D by notif
func getUpdateByNotif(notif smf_client.EventNotification) (bson.D, error) {
	var update bson.D
	var err error
	// TODO: implement other report types
	switch notif.GetEvent() {
	case smf_client.SMFEVENTANYOF_PDU_SES_EST:
		update, err = getUpdatePDU_SES_EST(notif)
	case smf_client.SMFEVENTANYOF_UE_IP_CH:
		update, err = getUpdateUE_IP_CH(notif)
	case smf_client.SMFEVENTANYOF_PLMN_CH:
		update, err = getUpdatePLMN_CH(notif)
	case smf_client.SMFEVENTANYOF_DDDS:
		update, err = getUpdateDDDS(notif)
	case smf_client.SMFEVENTANYOF_PDU_SES_REL:
		update, err = getUpdatePDU_SES_REL(notif)
	case smf_client.SMFEVENTANYOF_QOS_MON:
		update, err = getUpdateQOS_MON(notif)
	default:
		log.Printf("notif event %s is not supported currently",
			string(notif.GetEvent()))
		return nil, errors.New("invalid notif event")
	}
	if err != nil {
		return nil, err
	}
	return update, nil
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdatePDU_SES_EST - Create update bson.D in case of PDU SESS EST
func getUpdatePDU_SES_EST(notif smf_client.EventNotification) (bson.D, error) {
	adIpv4Addr, ok := notif.GetAdIpv4AddrOk()
	if !ok {
		return nil, errors.New("failed to get AdIpv4Addr")
	}
	dnn, ok := notif.GetDnnOk()
	if !ok {
		return nil, errors.New("failed to get Dnn")
	}
	pduSeId, ok := notif.GetPduSeIdOk()
	if !ok {
		return nil, errors.New("failed to get PduSeId")
	}
	pduSessType, ok := notif.GetPduSessTypeOk()
	if !ok {
		return nil, errors.New("failed to get PduSessType")
	}
	snssai, ok := notif.GetSnssaiOk()
	if !ok {
		return nil, errors.New("failed to get Snssai")
	}
	timeStamp := time.Now().Unix()
	push := pduSesEst{
		AdIpv4Addr:  adIpv4Addr,
		Dnn:         dnn,
		PduSeId:     pduSeId,
		PduSessType: pduSessType,
		Snssai:      snssai,
		TimeStamp:   timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"pdusesestlist": &push,
		}},
	}
	return update, nil
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdateUE_IP_CH - Create update bson.D in case of UE IP CH
func getUpdateUE_IP_CH(notif smf_client.EventNotification) (bson.D, error) {
	adIpv4Addr, ok := notif.GetAdIpv4AddrOk()
	if !ok {
		return nil, errors.New("failed to get AdIpv4Addr")
	}
	pduSeId, ok := notif.GetPduSeIdOk()
	if !ok {
		return nil, errors.New("failed to get PduSeId")
	}
	timeStamp := time.Now().Unix()
	push := ueIpCh{
		AdIpv4Addr: adIpv4Addr,
		PduSeId:    pduSeId,
		TimeStamp:  timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"ueipchlist": &push,
		}},
	}
	return update, nil
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdatePLMN_CH - Create update bson.D in case of PLMN CH
func getUpdatePLMN_CH(notif smf_client.EventNotification) (bson.D, error) {
	return nil, errors.New("getUpdatePLMN_CH not implemented")
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdateDDDS - Create update bson.D in case of DDDs
func getUpdateDDDS(notif smf_client.EventNotification) (bson.D, error) {
	dddStatus, ok := notif.GetDddStatusOk()
	if !ok {
		return nil, errors.New("failed to get DddStatus")
	}
	pduSeId, ok := notif.GetPduSeIdOk()
	if !ok {
		return nil, errors.New("failed to get PduSeId")
	}
	timeStamp := time.Now().Unix()
	push := ddds{DddStatus: dddStatus,
		PduSeId:   pduSeId,
		TimeStamp: timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"dddslist": &push,
		}},
	}
	return update, nil
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdatePDU_SES_REL - Create update bson.D in case of PDU SES REL
func getUpdatePDU_SES_REL(notif smf_client.EventNotification) (bson.D, error) {
	return nil, errors.New("getUpdatePDU_SES_REL not implemented")
}

// ----------------------------------------------------------------------------------------------------------------
// getUpdatePDU_SES_REL - Create update bson.D in case of QoS MON
func getUpdateQOS_MON(notif smf_client.EventNotification) (bson.D, error) {
	pduSeId, ok := notif.GetPduSeIdOk()
	if !ok {
		return nil, errors.New("failed to get PduSeId")
	}
	timeStamp := time.Now().Unix()
	// include "customized_data"
	push := qosMon{
		Customized_data: notif.CustomizedData,
		PduSeId:         pduSeId,
		TimeStamp:       timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"qosmonlist": &push,
		}},
	}
	return update, nil
}
