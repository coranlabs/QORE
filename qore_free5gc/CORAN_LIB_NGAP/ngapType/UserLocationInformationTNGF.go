package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UserLocationInformationTNGF struct {
	TNAPID       TNAPID
	IPAddress    TransportLayerAddress
	PortNumber   *PortNumber                                                  `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUserLocationInformationTNGFExtIEs `aper:"optional"`
}
