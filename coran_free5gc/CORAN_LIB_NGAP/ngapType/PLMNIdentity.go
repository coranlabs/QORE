package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PLMNIdentity struct {
	Value aper.OctetString `aper:"sizeLB:3,sizeUB:3"`
}
