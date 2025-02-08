package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	UPTransportLayerInformationPresentNothing int = iota /* No components present */
	UPTransportLayerInformationPresentGTPTunnel
	UPTransportLayerInformationPresentChoiceExtensions
)

type UPTransportLayerInformation struct {
	Present          int
	GTPTunnel        *GTPTunnel `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerUPTransportLayerInformationExtIEs
}
