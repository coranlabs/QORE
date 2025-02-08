package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	PresencePresentOptional    aper.Enumerated = 0
	PresencePresentConditional aper.Enumerated = 1
	PresencePresentMandatory   aper.Enumerated = 2
)

type Presence struct {
	Value aper.Enumerated `aper:"valueLB:0,valueUB:2"`
}
