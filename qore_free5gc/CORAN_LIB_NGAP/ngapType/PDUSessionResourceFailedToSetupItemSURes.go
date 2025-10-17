package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceFailedToSetupItemSURes struct {
	PDUSessionID                                PDUSessionID
	PDUSessionResourceSetupUnsuccessfulTransfer aper.OctetString
	IEExtensions                                *ProtocolExtensionContainerPDUSessionResourceFailedToSetupItemSUResExtIEs `aper:"optional"`
}
