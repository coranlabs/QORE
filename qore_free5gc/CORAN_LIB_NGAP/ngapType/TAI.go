package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type TAI struct {
	PLMNIdentity PLMNIdentity
	TAC          TAC
	IEExtensions *ProtocolExtensionContainerTAIExtIEs `aper:"optional"`
}
