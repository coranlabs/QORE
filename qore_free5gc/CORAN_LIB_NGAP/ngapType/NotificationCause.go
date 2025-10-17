package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	NotificationCausePresentFulfilled    aper.Enumerated = 0
	NotificationCausePresentNotFulfilled aper.Enumerated = 1
)

type NotificationCause struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
