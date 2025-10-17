package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type RATRestrictionsItem struct {
	PLMNIdentity              PLMNIdentity
	RATRestrictionInformation RATRestrictionInformation
	IEExtensions              *ProtocolExtensionContainerRATRestrictionsItemExtIEs `aper:"optional"`
}
