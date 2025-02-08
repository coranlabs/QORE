package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type SST struct {
	Value aper.OctetString `aper:"sizeLB:1,sizeUB:1"`
}
