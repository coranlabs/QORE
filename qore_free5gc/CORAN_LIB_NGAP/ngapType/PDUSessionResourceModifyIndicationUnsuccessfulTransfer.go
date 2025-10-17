package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceModifyIndicationUnsuccessfulTransfer struct {
	Cause        Cause                                                                                   `aper:"valueLB:0,valueUB:5"`
	IEExtensions *ProtocolExtensionContainerPDUSessionResourceModifyIndicationUnsuccessfulTransferExtIEs `aper:"optional"`
}
