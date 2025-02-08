package nas

import (
	"errors"
	"fmt"

	"github.com/coranlabs/CORAN_AMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
	"github.com/coranlabs/CORAN_AMF/Messages_handling_entity/gmm"
	nas "github.com/coranlabs/CORAN_LIB_NAS"

	//"github.com/coranlabs/CORAN_AMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/coranlabs/CORAN_LIB_UTIL/fsm"
)

func Dispatch(ue *context.AmfUe, accessType models.AccessType, procedureCode int64, msg *nas.Message) error {
	if msg.GmmMessage == nil {
		return errors.New("gmm Message is nil")
	}

	if msg.GsmMessage != nil {
		return errors.New("GSM Message should include in GMM Message")
	}

	if ue.State[accessType] == nil {
		return fmt.Errorf("UE State is empty (accessType=%q). Can't send GSM Message", accessType)
	}

	return gmm.GmmFSM.SendEvent(ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
		gmm.ArgAmfUe:         ue,
		gmm.ArgAccessType:    accessType,
		gmm.ArgNASMessage:    msg.GmmMessage,
		gmm.ArgProcedureCode: procedureCode,
	}, logger.GmmLog)
}
