package ngap

import (
	"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
)

func HandleMessage(ran *context.AmfRan, message *ngapType.NGAPPDU) {

	dispatchMain(ran, message)

}
