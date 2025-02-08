package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "gocoranlabs/lib/aper" if it uses "aper"

const (
	TNGFIDPresentNothing int = iota /* No components present */
	TNGFIDPresentTNGFID
	TNGFIDPresentChoiceExtensions
)

type TNGFID struct {
	Present          int             /* Choice Type */
	TNGFID           *aper.BitString `aper:"sizeLB:32,sizeUB:32"`
	ChoiceExtensions *ProtocolIESingleContainerTNGFIDExtIEs
}
