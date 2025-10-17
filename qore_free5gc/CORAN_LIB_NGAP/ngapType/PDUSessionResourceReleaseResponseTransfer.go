package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceReleaseResponseTransfer struct {
	IEExtensions *ProtocolExtensionContainerPDUSessionResourceReleaseResponseTransferExtIEs `aper:"optional"`
}
