package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceModifyItemModReq struct {
	PDUSessionID                            PDUSessionID
	NASPDU                                  *NASPDU `aper:"optional"`
	PDUSessionResourceModifyRequestTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceModifyItemModReqExtIEs `aper:"optional"`
}
