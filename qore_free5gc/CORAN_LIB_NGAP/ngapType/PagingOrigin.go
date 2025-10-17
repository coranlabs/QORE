package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	PagingOriginPresentNon3gpp aper.Enumerated = 0
)

type PagingOrigin struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
