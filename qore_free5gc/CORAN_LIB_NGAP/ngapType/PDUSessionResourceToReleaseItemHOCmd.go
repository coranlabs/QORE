package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceToReleaseItemHOCmd struct {
	PDUSessionID                            PDUSessionID
	HandoverPreparationUnsuccessfulTransfer aper.OctetString
	IEExtensions                            *ProtocolExtensionContainerPDUSessionResourceToReleaseItemHOCmdExtIEs `aper:"optional"`
}
