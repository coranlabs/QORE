package ngapType

// Need to import "gocoranlabs/lib/aper" if it uses "aper"

const (
	UserLocationInformationWAGFPresentNothing int = iota /* No components present */
	UserLocationInformationWAGFPresentGlobalLineID
	UserLocationInformationWAGFPresentHFCNodeID
	UserLocationInformationWAGFPresentChoiceExtensions
)

type UserLocationInformationWAGF struct {
	Present          int           /* Choice Type */
	GlobalLineID     *GlobalLineID `aper:"valueExt"`
	HFCNodeID        *HFCNodeID
	ChoiceExtensions *ProtocolIESingleContainerUserLocationInformationWAGFExtIEs
}
