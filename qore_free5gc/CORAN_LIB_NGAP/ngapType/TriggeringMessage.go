package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	TriggeringMessagePresentInitiatingMessage    aper.Enumerated = 0
	TriggeringMessagePresentSuccessfulOutcome    aper.Enumerated = 1
	TriggeringMessagePresentUnsuccessfullOutcome aper.Enumerated = 2
)

type TriggeringMessage struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}
