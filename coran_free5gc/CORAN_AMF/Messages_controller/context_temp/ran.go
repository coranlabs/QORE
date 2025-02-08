package context

import (
	"net"
	"sync"

	"github.com/coranlabs/CORAN_LIB_NGAP/ngapConvert"
	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	"github.com/sirupsen/logrus"
)

type Ran struct {
	RanPresent int
	RanId      *models.GlobalRanNodeId
	Name       string
	AnType     models.AccessType
	/* socket Connect*/
	Conn net.Conn
	/* Supported TA List */
	SupportedTAList []SupportedTAI

	/* RAN UE List */
	RanUeList sync.Map // RanUeNgapId as key

	/* logger */
	Log *logrus.Entry
}

func (ran *Ran) SetRanId(ranNodeId *ngapType.GlobalRANNodeID) {
	ranId := ngapConvert.RanIdToModels(*ranNodeId)
	ran.RanPresent = ranNodeId.Present
	ran.RanId = &ranId
	if ranNodeId.Present == ngapType.GlobalRANNodeIDPresentGlobalN3IWFID ||
		ranNodeId.Present == ngapType.GlobalRANNodeIDPresentChoiceExtensions {
		ran.AnType = models.AccessType_NON_3_GPP_ACCESS
	} else {
		ran.AnType = models.AccessType__3_GPP_ACCESS
	}
}
