package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type SecurityContext struct {
	NextHopChainingCount NextHopChainingCount
	NextHopNH            SecurityKey
	IEExtensions         *ProtocolExtensionContainerSecurityContextExtIEs `aper:"optional"`
}
