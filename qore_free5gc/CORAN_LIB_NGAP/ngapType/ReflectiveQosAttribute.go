package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	ReflectiveQosAttributePresentSubjectTo aper.Enumerated = 0
)

type ReflectiveQosAttribute struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
