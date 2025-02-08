package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type EPSTAC struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}
