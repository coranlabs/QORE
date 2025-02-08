package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	HandoverTypePresentIntra5gs    aper.Enumerated = 0
	HandoverTypePresentFivegsToEps aper.Enumerated = 1
	HandoverTypePresentEpsTo5gs    aper.Enumerated = 2
)

type HandoverType struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
}
