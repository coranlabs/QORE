package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type RepetitionPeriod struct {
	Value int64 `aper:"valueLB:0,valueUB:131071"`
}
