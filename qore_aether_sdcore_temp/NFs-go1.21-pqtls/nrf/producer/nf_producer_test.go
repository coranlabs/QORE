// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package producer_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/omec-project/nrf/dbadapter"
	"github.com/omec-project/nrf/factory"
	"github.com/omec-project/nrf/producer"
	"github.com/omec-project/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

type MockMongoDBClient struct {
	dbadapter.DBInterface
}

func init() {
	factory.InitConfigFactory("../nrfTest/nrfcfg.yaml")
}

func (db *MockMongoDBClient) RestfulAPIGetOne(collName string, filter bson.M) map[string]interface{} {
	fmt.Println("Called Mock RestfulAPIGetOne")
	return nil
}

func (db *MockMongoDBClient) RestfulAPIGetMany(collName string, filter bson.M) []map[string]interface{} {
	fmt.Println("Called Mock RestfulAPIGetMany")
	return nil
}
func (db *MockMongoDBClient) PutOneWithTimeout(collName string, filter bson.M, putData map[string]interface{}, timeout int32, timeField string) bool {
	fmt.Println("Called Mock PutOneWithTimeout")
	return true
}
func (db *MockMongoDBClient) RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) bool {
	fmt.Println("Called Mock RestfulAPIPutOne")
	return true
}
func (db *MockMongoDBClient) RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) bool {
	fmt.Println("Called Mock RestfulAPIPutOneNotUpdate")
	return true
}
func (db *MockMongoDBClient) RestfulAPIPutMany(collName string, filterArray []bson.M, putDataArray []map[string]interface{}) bool {
	fmt.Println("Called Mock RestfulAPIPutMany")
	return true
}
func (db *MockMongoDBClient) RestfulAPIDeleteOne(collName string, filter bson.M) {
	fmt.Println("Called Mock RestfulAPIDeleteOne")
}
func (db *MockMongoDBClient) RestfulAPIDeleteMany(collName string, filter bson.M) {
	fmt.Println("Called Mock RestfulAPIDeleteMany")
}
func (db *MockMongoDBClient) RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) bool {
	fmt.Println("Called Mock RestfulAPIMergePatch")
	return true
}
func (db *MockMongoDBClient) RestfulAPIJSONPatch(collName string, filter bson.M, patchJSON []byte) bool {
	fmt.Println("Called Mock RestfulAPIJSONPatch")
	return true
}
func (db *MockMongoDBClient) RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchJSON []byte, dataName string) bool {
	fmt.Println("Called Mock RestfulAPIJSONPatchExtend")
	return true
}
func (db *MockMongoDBClient) RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) bool {
	fmt.Println("Called Mock RestfulAPIPost")
	return true
}
func (db *MockMongoDBClient) RestfulAPIPostMany(collName string, filter bson.M, postDataArray []interface{}) bool {
	fmt.Println("Called Mock RestfulAPIPost")
	return true
}

func TestNFRegisterProcedure(t *testing.T) {
	dbadapter.DBClient = &MockMongoDBClient{}
	var nf models.NfProfile
	nf.NfType = models.NfType_PCF
	nf.NfInstanceId = uuid.New().String()
	nf.NfStatus = models.NfStatus_REGISTERED
	_, _, err := producer.NFRegisterProcedure(nf)
	if err != nil {
		t.Errorf("testcase failed: %v", err)
	}
}
