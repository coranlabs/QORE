package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	DirectForwardingPathAvailabilityPresentDirectPathAvailable aper.Enumerated = 0
)

type DirectForwardingPathAvailability struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
