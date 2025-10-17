package mongodb

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"

	// "github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	// "github.com/coranlabs/CORAN_UDR/Application_entity/server/util"
	// "github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
	openapi "github.com/coranlabs/CORAN_LIB_OPENAPI"
	"github.com/coranlabs/CORAN_LIB_UTIL/mongoapi"
	"github.com/coranlabs/CORAN_UDR/Application_entity/logger"
	"github.com/coranlabs/CORAN_UDR/Application_entity/pkg/factory"
	"github.com/coranlabs/CORAN_UDR/Application_entity/server/util"
)

type MongoDbConnector struct {
	*factory.Mongodb
}

func NewMongoDbConnector(mongo *factory.Mongodb) MongoDbConnector {
	return MongoDbConnector{
		Mongodb: mongo,
	}
}

func (m MongoDbConnector) PatchDataToDBAndNotify(
	collName string, ueId string, patchItem []models.PatchItem, filter bson.M,
) (origValue, newValue map[string]interface{}, err error) {
	origValue, err = mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		return
	}

	patchJSON, err := json.Marshal(patchItem)
	if err != nil {
		return
	}

	if err = mongoapi.RestfulAPIJSONPatch(collName, filter, patchJSON); err != nil {
		return
	}

	newValue, err = mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		return
	}

	return
}

func (m MongoDbConnector) GetDataFromDB(
	collName string, filter bson.M) (
	map[string]interface{}, *models.ProblemDetails,
) {
	data, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		return nil, openapi.ProblemDetailsSystemFailure(err.Error())
	}
	if data == nil {
		return nil, util.ProblemDetailsNotFound("DATA_NOT_FOUND")
	}
	return data, nil
}

func (m MongoDbConnector) GetDataFromDBWithArg(collName string, filter bson.M, strength int) (
	map[string]interface{}, *models.ProblemDetails,
) {
	data, err := mongoapi.RestfulAPIGetOne(collName, filter, strength)
	if err != nil {
		return nil, openapi.ProblemDetailsSystemFailure(err.Error())
	}
	if data == nil {
		logger.ConsumerLog.Errorln("filter: ", filter)
		return nil, util.ProblemDetailsNotFound("DATA_NOT_FOUND")
	}

	return data, nil
}

func (m MongoDbConnector) DeleteDataFromDB(collName string, filter bson.M) {
	if err := mongoapi.RestfulAPIDeleteOne(collName, filter); err != nil {
		logger.DataRepoLog.Errorf("deleteDataFromDB: %+v", err)
	}
}
