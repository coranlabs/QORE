package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type SliceOverloadItem struct {
	SNSSAI       SNSSAI                                             `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerSliceOverloadItemExtIEs `aper:"optional"`
}
