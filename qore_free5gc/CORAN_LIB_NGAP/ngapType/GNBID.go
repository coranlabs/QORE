package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	GNBIDPresentNothing int = iota /* No components present */
	GNBIDPresentGNBID
	GNBIDPresentChoiceExtensions
)

type GNBID struct {
	Present          int
	GNBID            *aper.BitString `aper:"sizeLB:22,sizeUB:32"`
	ChoiceExtensions *ProtocolIESingleContainerGNBIDExtIEs
}
