package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	UENGAPIDsPresentNothing int = iota /* No components present */
	UENGAPIDsPresentUENGAPIDPair
	UENGAPIDsPresentAMFUENGAPID
	UENGAPIDsPresentChoiceExtensions
)

type UENGAPIDs struct {
	Present          int
	UENGAPIDPair     *UENGAPIDPair `aper:"valueExt"`
	AMFUENGAPID      *AMFUENGAPID
	ChoiceExtensions *ProtocolIESingleContainerUENGAPIDsExtIEs
}
