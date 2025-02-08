package Messages_controller

import (
	"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
	"github.com/coranlabs/CORAN_AMF/Messages_handling_entity/ngap"

	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
)

func HandleNgap(ran *context.AmfRan, message *ngapType.NGAPPDU) {

	ngap.HandleMessage(ran, message)

}
