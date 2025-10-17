package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "gocoranlabs/lib/aper" if it uses "aper"

const (
	WAGFIDPresentNothing int = iota /* No components present */
	WAGFIDPresentWAGFID
	WAGFIDPresentChoiceExtensions
)

type WAGFID struct {
	Present          int             /* Choice Type */
	WAGFID           *aper.BitString `aper:"sizeExt,sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerWAGFIDExtIEs
}
