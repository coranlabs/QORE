package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	NgENBIDPresentNothing int = iota /* No components present */
	NgENBIDPresentMacroNgENBID
	NgENBIDPresentShortMacroNgENBID
	NgENBIDPresentLongMacroNgENBID
	NgENBIDPresentChoiceExtensions
)

type NgENBID struct {
	Present           int
	MacroNgENBID      *aper.BitString `aper:"sizeLB:20,sizeUB:20"`
	ShortMacroNgENBID *aper.BitString `aper:"sizeLB:18,sizeUB:18"`
	LongMacroNgENBID  *aper.BitString `aper:"sizeLB:21,sizeUB:21"`
	ChoiceExtensions  *ProtocolIESingleContainerNgENBIDExtIEs
}
