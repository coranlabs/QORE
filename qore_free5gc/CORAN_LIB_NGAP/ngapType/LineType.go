package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "gocoranlabs/lib/aper" if it uses "aper"

const ( /* Enum Type */
	LineTypePresentDsl aper.Enumerated = 0
	LineTypePresentPon aper.Enumerated = 1
)

type LineType struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
