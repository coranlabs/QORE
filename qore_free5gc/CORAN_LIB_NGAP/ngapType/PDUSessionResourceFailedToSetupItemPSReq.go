package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceFailedToSetupItemPSReq struct {
	PDUSessionID                         PDUSessionID
	PathSwitchRequestSetupFailedTransfer aper.OctetString
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemPSReqExtIEs `aper:"optional"`
}
