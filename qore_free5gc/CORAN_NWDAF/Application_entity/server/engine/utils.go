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

package engine

import (
	"context"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------------------------------------------------------------
// InitConfig - Initialize global variables (cfg and mongoClient)
func InitConfig() {

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientOptions := options.Client().ApplyURI(config.Database.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to MongoDB.")
	mongoClient = client
}

// ------------------------------------------------------------------------------
// getAnaReqValues - Get Analytics-request start-ts and ent-ts values
func getExtraReportReq(engineReqData EngineReqData) (int64, int64) {
	startTs, endTs := int64(0), int64(0)
	if !engineReqData.StartTs.IsZero() {
		startTs = engineReqData.StartTs.Unix()
	}
	if !engineReqData.EndTs.IsZero() {
		endTs = engineReqData.EndTs.Unix()
	}
	return startTs, endTs
}

// ------------------------------------------------------------------------------
// getTimeStampCondition - Get timestamp confition to be used by filters
func getTimeStampCondition(startTs int64, endTs int64) bson.D {
	timeStamp := calculateTimeStamp(startTs, endTs)
	switch timeStamp {
	case 0:
		return bson.D{{"$gt", startTs}, {"$lt", endTs}}
	default:
		return bson.D{{"$lt", timeStamp}}
	}
}

// ------------------------------------------------------------------------------
// calculateTimeStamp - Set timestamp according to cases 1-2-3 presented below
/*
	case 1. No startTs or endTs -> number of documents before NOW
	case 2A. startTs present, no endTs -> number of documents before startTs
	case 2B. endTs present, no startTs -> number of documents before endTs
	case 3. Both startTs and endTs present -> number of documents during period
*/
func calculateTimeStamp(startTs int64, endTs int64) int64 {
	timeStamp := int64(0)
	if endTs == 0 {
		if startTs == 0 {
			log.Printf("Case 1. No startTs or endTs-> documents before NOW")
			timeStamp = time.Now().Unix()
		} else {
			log.Printf("Case 2A. startTs present, no endTs -> documents before startTs")
			timeStamp = startTs
		}
	} else {
		if startTs == 0 {
			log.Printf("Case 2B. no startTs, endTs present -> documents before endTs")
			timeStamp = endTs
		} else {
			log.Printf("Case 3. Both startTs and endTs present -> documents during period")
		}
	}
	return timeStamp
}

// ------------------------------------------------------------------------------
// MatchTimeStamp - Returns true if notifTimeStamp corresponds to what we are searching for.
func matchTimeStamp(notifTimeStamp int64, timeStamp int64, startTs int64, endTs int64) bool {
	if timeStamp == 0 {
		return (notifTimeStamp > startTs && notifTimeStamp < endTs)
	} else {
		return notifTimeStamp < timeStamp
	}
}
