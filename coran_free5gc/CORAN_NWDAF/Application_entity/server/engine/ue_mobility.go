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
 * Description: This file contains functions related to UE mobitlity event ID.
 */

package engine

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------------------------------------------------------------------------------
// ueMob - get ue mobility information.
func ueMob(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		log.Printf("Getting Ue Mobility from DB")
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
		// Construct the filter based on SUPI and timestamp interval
		filter := GetFilterUeMob(engineReqData)
		// Retrieve the documents matching the filter
		db := mongoClient.Database(config.Database.DbName)
		collection := db.Collection(config.Database.CollectionAmfName)
		// Find the IMSI document
		var result bson.M
		var loc []bson.M
		err = collection.FindOne(context.Background(), filter).Decode(&result)
		if err == nil {
			// Extract the location list from the result
			locationList := result["locationlist"].(primitive.A)
			// Variables to store the last user location and its timestamp
			var lastUserLocation bson.M
			var lastTimestamp int64
			// get startTs and endTs
			startTs, endTs := getExtraReportReq(engineReqData)
			// Iterate over the location list and extract the nrlocation
			for _, location := range locationList {
				locationMap := location.(primitive.M)
				userLocation := locationMap["userlocation"].(bson.M)
				timestamp := locationMap["timestamp"].(int64)
				// Check if either startTs or endTs is 0
				if startTs == 0 || endTs == 0 {
					// Update the last user location and its timestamp
					if timestamp > lastTimestamp {
						lastUserLocation = userLocation
						lastTimestamp = timestamp
					}
				} else {
					// Check if the timestamp falls within the desired range
					if timestamp > startTs && timestamp < endTs {
						loc = append(loc, userLocation)
					}
				}
			}
			// If either startTs or endTs is 0, append the last user location
			if startTs == 0 || endTs == 0 {
				loc = append(loc, lastUserLocation)
			}
		}
		response := bson.M{
			"loc": loc,
		}
		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(response)
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
// GetFilterUeMob - Get filter to retrieve documents according to startTs and endTs values
func GetFilterUeMob(engineReqData EngineReqData) bson.M {
	log.Printf("Constructing Ue Mobility filter for DB search ...")
	return bson.M{"_id": engineReqData.Supi}
}
