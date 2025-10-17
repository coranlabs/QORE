package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	CPTransportLayerInformationPresentNothing int = iota /* No components present */
	CPTransportLayerInformationPresentEndpointIPAddress
	CPTransportLayerInformationPresentChoiceExtensions
)

type CPTransportLayerInformation struct {
	Present           int
	EndpointIPAddress *TransportLayerAddress
	ChoiceExtensions  *ProtocolIESingleContainerCPTransportLayerInformationExtIEs
}
