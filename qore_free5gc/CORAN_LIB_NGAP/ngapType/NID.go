package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

type NID struct {
	Value aper.BitString `aper:"sizeLB:44,sizeUB:44"`
}
