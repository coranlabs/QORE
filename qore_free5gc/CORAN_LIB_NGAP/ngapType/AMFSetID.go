package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type AMFSetID struct {
	Value aper.BitString `aper:"sizeLB:10,sizeUB:10"`
}
