package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	TimerApproachForGUAMIRemovalPresentApplyTimer aper.Enumerated = 0
)

type TimerApproachForGUAMIRemoval struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
