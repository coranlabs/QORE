package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UserLocationInformationNR struct {
	NRCGI        NRCGI                                                      `aper:"valueExt"`
	TAI          TAI                                                        `aper:"valueExt"`
	TimeStamp    *TimeStamp                                                 `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUserLocationInformationNRExtIEs `aper:"optional"`
}
