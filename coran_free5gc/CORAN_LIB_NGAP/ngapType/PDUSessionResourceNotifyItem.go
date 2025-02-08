package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceNotifyItem struct {
	PDUSessionID                     PDUSessionID
	PDUSessionResourceNotifyTransfer aper.OctetString
	IEExtensions                     *ProtocolExtensionContainerPDUSessionResourceNotifyItemExtIEs `aper:"optional"`
}
