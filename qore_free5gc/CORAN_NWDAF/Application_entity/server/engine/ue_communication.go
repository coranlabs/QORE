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
 * Description: This file contains functions related to UE communication event ID.
 */

package engine

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------------------------------------------------------------------------------
// ueComm - get ue communications statistics.
func ueComm(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		log.Printf("Getting Ue Communication from DB")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var engineReqData EngineReqData
		err = json.Unmarshal(body, &engineReqData)
		if err != nil {
			http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
			return
		}
		// get filter to retrieve documents according to startTs and endTs values
		filter := GetFilterUeComm(engineReqData)
		// get collection from database
		db := mongoClient.Database(config.Database.DbName)
		collection := db.Collection(config.Database.CollectionSmfName)
		// Retrieve Documents from DB
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			http.Error(w, "Error finding documents", http.StatusInternalServerError)
			return
		}
		commDur := int32(0)
		ulVolSlice, dlVolSlice := make([]int64, 0), make([]int64, 0)
		startTs, endTs := getExtraReportReq(engineReqData)
		timeStamp := calculateTimeStamp(startTs, endTs)
		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				http.Error(w, "Error decoding document", http.StatusInternalServerError)
				return
			}
			qosMonList, ok := result["qosmonlist"].(primitive.A)
			if !ok {
				http.Error(w, "Invalid qosmonlist type in document", http.StatusInternalServerError)
				return
			}
			for _, qosMonElem := range qosMonList {
				qosMonMap, ok := qosMonElem.(bson.M)
				if !ok {
					http.Error(w, "Invalid qosMonElem type in document", http.StatusInternalServerError)
					return
				}
				qosTimestamp := qosMonMap["timestamp"].(int64)
				if matchTimeStamp(qosTimestamp, timeStamp, startTs, endTs) {
					usageReport := qosMonMap["customized_data"].(bson.M)["usagereport"].(bson.M)
					volume := usageReport["volume"].(bson.M)
					duration := usageReport["duration"].(int32)
					ulVolSlice = append(ulVolSlice, volume["uplink"].(int64))
					dlVolSlice = append(dlVolSlice, volume["downlink"].(int64))
					commDur += duration
				}
			}
		}
		ulVol, ulVolVariance := calculateSumAndVariance(ulVolSlice)
		dlVol, dlVolVariance := calculateSumAndVariance(dlVolSlice)
		ueCommResp := UeCommResp{
			CommDur:       commDur,
			Ts:            engineReqData.StartTs,
			UlVol:         ulVol,
			UlVolVariance: ulVolVariance,
			DlVol:         dlVol,
			DlVolVariance: dlVolVariance,
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(ueCommResp)
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------
// GetFilterUeComm - Get filter to retrieve documents according to startTs and endTs values
func GetFilterUeComm(engineReqData EngineReqData) bson.D {
	log.Printf("Constructing Ue Communication filter for DB search ...")
	// get startTs and endTs
	startTs, endTs := getExtraReportReq(engineReqData)
	// get timeStamp condition to inject in mongo filter
	timeStampCondition := getTimeStampCondition(startTs, endTs)
	return bson.D{{"qosmonlist",
		bson.M{"$elemMatch": bson.M{
			"timestamp": timeStampCondition,
		}},
	}}
}

// -----------------------------------------------------------------------------
// calculateSumAndVariance - Calculates sum and variance
func calculateSumAndVariance(volSlice []int64) (int64, float32) {

	if len(volSlice) == 0 {
		return int64(0), float32(0)
	}
	volSum := int64(0)
	volVariance := float64(0)
	for _, vol := range volSlice {
		volSum += vol
	}
	// sum the square of the mean subtracted from each element
	volMean := float64(volSum) / float64(len(volSlice))
	for _, vol := range volSlice {
		volVariance += (float64(vol) - volMean) * (float64(vol) - volMean)
	}
	// divide variance by the slice length and take square root
	volVarianceResult := float32(math.Sqrt(volVariance / float64(len(volSlice))))
	return volSum, volVarianceResult
}
