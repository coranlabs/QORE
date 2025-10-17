package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type TargeteNBID struct {
	GlobalENBID    GlobalNgENBID                                `aper:"valueExt"`
	SelectedEPSTAI EPSTAI                                       `aper:"valueExt"`
	IEExtensions   *ProtocolExtensionContainerTargeteNBIDExtIEs `aper:"optional"`
}
