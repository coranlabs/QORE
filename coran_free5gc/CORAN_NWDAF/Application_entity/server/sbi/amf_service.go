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
 * Description: This file contains functions related amf post notifications.
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

	amf_client "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/clients/amfclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------------------------------------------------------------
func storeAmfNotificationOnDB(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST":
		log.Printf("Storing AMF notification in Database")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var amfNotification *amf_client.AmfEventNotification
		err = json.Unmarshal(body, &amfNotification)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}
		reportList, ok := amfNotification.GetReportListOk()
		if !ok {
			http.Error(w, "Error in getting ReportList from AMF notification", http.StatusBadRequest)
			return
		}
		databaseName := config.Database.DbName
		collectionName := config.Database.CollectionAmfName
		amfCollection := mongoClient.Database(databaseName).Collection(collectionName)
		opts := options.Update().SetUpsert(true)
		// store reports one by one
		for _, report := range reportList {
			oid := report.GetSupi()
			if oid == "" {
				http.Error(w, "supi not found in report, cannot create object id", http.StatusBadRequest)
				return
			}
			update, err := getUpdateByReport(report)
			if err != nil {
				http.Error(w, "error in getUpdateByReport", http.StatusBadRequest)
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			res, err := amfCollection.UpdateByID(ctx, oid, update, opts)
			if err != nil {
				http.Error(w, "error in updating the AMF collection", http.StatusBadRequest)
				return
			}
			if res.MatchedCount != 0 {
				log.Printf("Matched and updated an existing notification report from Amf")
			}
			if res.UpsertedCount != 0 {
				log.Printf("Inserted a new notification report from Amf with ID %v\n",
					res.UpsertedID)
			}
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------
// getUpdateByReport - Return update bson.D by report
func getUpdateByReport(report amf_client.AmfEventReport) (bson.D, error) {
	var update bson.D
	var err error
	switch report.GetType() {
	case amf_client.AMFEVENTTYPEANYOF_REGISTRATION_STATE_REPORT:
		update, err = getUpdateRegistration(report)
	case amf_client.AMFEVENTTYPEANYOF_LOCATION_REPORT:
		update, err = getUpdateLocation(report)
	case amf_client.AMFEVENTTYPEANYOF_LOSS_OF_CONNECTIVITY:
		update, err = getUpdateLossOfConnectivity(report)
	default:
		log.Printf("report type %s is not supported currently", string(report.GetType()))
		return nil, errors.New("invalid report type")
	}
	if err != nil {
		return nil, err
	}
	return update, nil
}

// ------------------------------------------------------------------------------
// getUpdateRegistration - Create update bson.D in case of registration
func getUpdateRegistration(report amf_client.AmfEventReport) (bson.D, error) {
	rmInfoList, ok := report.GetRmInfoListOk()
	if !ok {
		return nil, errors.New("failed to get RmInfoList")
	}
	timeStamp := time.Now().Unix()
	// TODO: fix (get rid of) the "RmStateAnyOf" field
	push := rmInfo{
		RmInfo:    rmInfoList[len(rmInfoList)-1],
		TimeStamp: timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"rminfolist": &push,
		}},
	}
	return update, nil
}

// ------------------------------------------------------------------------------
// getUpdateLocation - Create update bson.D in case of Location
func getUpdateLocation(report amf_client.AmfEventReport) (bson.D, error) {
	locationObj, ok := report.GetLocationOk()
	if !ok {
		return nil, errors.New("failed to get Location")
	}
	timeStamp := time.Now().Unix()
	push := location{UserLocation: *locationObj, TimeStamp: timeStamp}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"locationlist": &push,
		}},
	}
	return update, nil
}

// ------------------------------------------------------------------------------
// getUpdateLossOfConnectivity - Create update bson.D in case of Loss of connectivity
func getUpdateLossOfConnectivity(report amf_client.AmfEventReport) (bson.D, error) {
	lossOfConnectReasonObj, ok := report.GetLossOfConnectReasonOk()
	if !ok {
		return nil, errors.New("failed to get lossOfConnectReason")
	}
	timeStamp := time.Now().Unix()
	push := lossOfConnectReason{
		LossOfConnectReason: *lossOfConnectReasonObj.LossOfConnectivityReasonAnyOf,
		TimeStamp:           timeStamp,
	}
	update := bson.D{
		{"$set", bson.D{
			{"lastmodified", timeStamp},
		}},
		{"$push", bson.M{
			"lossofconnectreasonlist": &push,
		}},
	}
	return update, nil
}
