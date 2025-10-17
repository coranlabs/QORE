package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UPTransportLayerInformationItem struct {
	NGUUPTNLInformation UPTransportLayerInformation                                      `aper:"valueLB:0,valueUB:1"`
	IEExtensions        *ProtocolExtensionContainerUPTransportLayerInformationItemExtIEs `aper:"optional"`
}
