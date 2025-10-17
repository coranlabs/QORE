package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type AMFPointer struct {
	Value aper.BitString `aper:"sizeLB:6,sizeUB:6"`
}
