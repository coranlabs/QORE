package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UPTransportLayerInformationPairItem struct {
	ULNGUUPTNLInformation UPTransportLayerInformation                                          `aper:"valueLB:0,valueUB:1"`
	DLNGUUPTNLInformation UPTransportLayerInformation                                          `aper:"valueLB:0,valueUB:1"`
	IEExtensions          *ProtocolExtensionContainerUPTransportLayerInformationPairItemExtIEs `aper:"optional"`
}
