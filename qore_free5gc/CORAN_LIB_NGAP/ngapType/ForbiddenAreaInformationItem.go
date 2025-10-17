package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type ForbiddenAreaInformationItem struct {
	PLMNIdentity  PLMNIdentity
	ForbiddenTACs ForbiddenTACs
	IEExtensions  *ProtocolExtensionContainerForbiddenAreaInformationItemExtIEs `aper:"optional"`
}
