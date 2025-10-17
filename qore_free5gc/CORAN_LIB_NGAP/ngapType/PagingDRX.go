package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	PagingDRXPresentV32  aper.Enumerated = 0
	PagingDRXPresentV64  aper.Enumerated = 1
	PagingDRXPresentV128 aper.Enumerated = 2
	PagingDRXPresentV256 aper.Enumerated = 3
)

type PagingDRX struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:3"`
}
