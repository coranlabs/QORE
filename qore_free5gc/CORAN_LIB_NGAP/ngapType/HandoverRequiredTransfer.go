package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type HandoverRequiredTransfer struct {
	DirectForwardingPathAvailability *DirectForwardingPathAvailability                         `aper:"optional"`
	IEExtensions                     *ProtocolExtensionContainerHandoverRequiredTransferExtIEs `aper:"optional"`
}
