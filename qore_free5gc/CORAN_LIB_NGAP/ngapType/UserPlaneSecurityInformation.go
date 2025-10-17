package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UserPlaneSecurityInformation struct {
	SecurityResult     SecurityResult                                                `aper:"valueExt"`
	SecurityIndication SecurityIndication                                            `aper:"valueExt"`
	IEExtensions       *ProtocolExtensionContainerUserPlaneSecurityInformationExtIEs `aper:"optional"`
}
