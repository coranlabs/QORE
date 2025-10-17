package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	ExpectedHOIntervalPresentSec15    aper.Enumerated = 0
	ExpectedHOIntervalPresentSec30    aper.Enumerated = 1
	ExpectedHOIntervalPresentSec60    aper.Enumerated = 2
	ExpectedHOIntervalPresentSec90    aper.Enumerated = 3
	ExpectedHOIntervalPresentSec120   aper.Enumerated = 4
	ExpectedHOIntervalPresentSec180   aper.Enumerated = 5
	ExpectedHOIntervalPresentLongTime aper.Enumerated = 6
)

type ExpectedHOInterval struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:6"`
}
