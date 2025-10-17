package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type SNSSAI struct {
	SST          SST
	SD           *SD                                     `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerSNSSAIExtIEs `aper:"optional"`
}
