package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	ConfidentialityProtectionResultPresentPerformed    aper.Enumerated = 0
	ConfidentialityProtectionResultPresentNotPerformed aper.Enumerated = 1
)

type ConfidentialityProtectionResult struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
