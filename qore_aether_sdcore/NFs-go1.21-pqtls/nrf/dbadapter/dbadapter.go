// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
package dbadapter

import (
	"context"
	"log"

	"github.com/omec-project/util/mongoapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBInterface interface {
	RestfulAPIGetOne(collName string, filter bson.M) (map[string]interface{}, error)
	RestfulAPIGetMany(collName string, filter bson.M) ([]map[string]interface{}, error)
	//	PutOneWithTimeout(collName string, filter bson.M, putData map[string]interface{}, timeout int32, timeField string) bool
	RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) (bool, error)
	RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) (bool, error)
	RestfulAPIDeleteOne(collName string, filter bson.M) error
	RestfulAPIDeleteMany(collName string, filter bson.M) error
	RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) error
	RestfulAPIJSONPatch(collName string, filter bson.M, patchJSON []byte) error
	RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchJSON []byte, dataName string) error
	RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) (bool, error)
	RestfulAPIPutMany(collName string, filterArray []primitive.M, putDataArray []map[string]interface{}) error
}

var DBClient DBInterface = nil

type MongoDBClient struct {
	mongoapi.MongoClient
}

func iterateChangeStream(routineCtx context.Context, stream *mongo.ChangeStream) {
	log.Println("iterate change stream for timeout")
	defer stream.Close(routineCtx)
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		log.Println("iterate stream : ", data)
	}
}

func ConnectToDBClient(dbName string, url string, enableStream bool, nfProfileExpiryEnable bool) DBInterface {
	for {
		MongoClient, _ := mongoapi.NewMongoClient(url, dbName)
		if MongoClient != nil {
			log.Println("MongoDB Connection Successful")
			DBClient = MongoClient
			break
		} else {
			log.Println("MongoDB Connection Failed")
		}
	}

	db := DBClient.(*mongoapi.MongoClient)
	if enableStream {
		log.Println("MongoDB Change stream Enabled")
		database := db.Client.Database(dbName)
		NfProfileColl := database.Collection("NfProfile")
		//create stream to monitor actions on the collection
		NfProfStream, err := NfProfileColl.Watch(context.TODO(), mongo.Pipeline{})
		if err != nil {
			panic(err)
		}
		routineCtx, _ := context.WithCancel(context.Background())
		//run routine to get messages from stream
		go iterateChangeStream(routineCtx, NfProfStream)
	}

	if nfProfileExpiryEnable {
		log.Println("NfProfile document expiry enabled")
		ret := db.RestfulAPICreateTTLIndex("NfProfile", 0, "expireAt")
		if ret {
			log.Println("TTL Index created for Field : expireAt in Collection : NfProfile")
		} else {
			log.Println("TTL Index exists for Field : expireAt in Collection : NfProfile")
		}
	}
	return DBClient
}

func (db *MongoDBClient) RestfulAPIGetOne(collName string, filter bson.M) (map[string]interface{}, error) {
	return db.RestfulAPIGetOne(collName, filter)
}

func (db *MongoDBClient) RestfulAPIGetMany(collName string, filter bson.M) []map[string]interface{} {
	return db.RestfulAPIGetMany(collName, filter)
}
func (db *MongoDBClient) PutOneWithTimeout(collName string, filter bson.M, putData map[string]interface{}, timeout int32, timeField string) bool {
	return db.RestfulAPIPutOneTimeout(collName, filter, putData, timeout, timeField)
}
func (db *MongoDBClient) RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) bool {
	return db.RestfulAPIPutOne(collName, filter, putData)
}
func (db *MongoDBClient) RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) bool {
	return db.RestfulAPIPutOneNotUpdate(collName, filter, putData)
}
func (db *MongoDBClient) RestfulAPIPutMany(collName string, filterArray []primitive.M, putDataArray []map[string]interface{}) bool {
	return db.RestfulAPIPutMany(collName, filterArray, putDataArray)
}
func (db *MongoDBClient) RestfulAPIDeleteOne(collName string, filter bson.M) {
	db.RestfulAPIDeleteOne(collName, filter)
}
func (db *MongoDBClient) RestfulAPIDeleteMany(collName string, filter bson.M) {
	db.RestfulAPIDeleteMany(collName, filter)
}
func (db *MongoDBClient) RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) bool {
	return db.RestfulAPIMergePatch(collName, filter, patchData)
}
func (db *MongoDBClient) RestfulAPIJSONPatch(collName string, filter bson.M, patchJSON []byte) bool {
	return db.RestfulAPIJSONPatch(collName, filter, patchJSON)
}
func (db *MongoDBClient) RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchJSON []byte, dataName string) bool {
	return db.RestfulAPIJSONPatchExtend(collName, filter, patchJSON, dataName)
}
func (db *MongoDBClient) RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) bool {
	return db.RestfulAPIPost(collName, filter, postData)
}
func (db *MongoDBClient) RestfulAPIPostMany(collName string, filter bson.M, postDataArray []interface{}) bool {
	return db.RestfulAPIPostMany(collName, filter, postDataArray)
}
