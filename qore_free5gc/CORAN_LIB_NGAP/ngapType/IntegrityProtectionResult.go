package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	IntegrityProtectionResultPresentPerformed    aper.Enumerated = 0
	IntegrityProtectionResultPresentNotPerformed aper.Enumerated = 1
)

type IntegrityProtectionResult struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
