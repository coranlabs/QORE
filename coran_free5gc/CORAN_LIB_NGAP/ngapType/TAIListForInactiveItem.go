package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type TAIListForInactiveItem struct {
	TAI          TAI                                                     `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerTAIListForInactiveItemExtIEs `aper:"optional"`
}
