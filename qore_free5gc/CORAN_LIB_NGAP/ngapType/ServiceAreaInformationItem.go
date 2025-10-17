package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type ServiceAreaInformationItem struct {
	PLMNIdentity   PLMNIdentity
	AllowedTACs    *AllowedTACs                                                `aper:"optional"`
	NotAllowedTACs *NotAllowedTACs                                             `aper:"optional"`
	IEExtensions   *ProtocolExtensionContainerServiceAreaInformationItemExtIEs `aper:"optional"`
}
