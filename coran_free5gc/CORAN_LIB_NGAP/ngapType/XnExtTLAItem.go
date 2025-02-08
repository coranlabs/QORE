package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type XnExtTLAItem struct {
	IPsecTLA     *TransportLayerAddress                        `aper:"optional"`
	GTPTLAs      *XnGTPTLAs                                    `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerXnExtTLAItemExtIEs `aper:"optional"`
}
