package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type ExpectedUEMovingTrajectoryItem struct {
	NGRANCGI         NGRANCGI                                                        `aper:"valueLB:0,valueUB:2"`
	TimeStayedInCell *int64                                                          `aper:"valueLB:0,valueUB:4095,optional"`
	IEExtensions     *ProtocolExtensionContainerExpectedUEMovingTrajectoryItemExtIEs `aper:"optional"`
}
