package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type QosFlowAcceptedItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowAcceptedItemExtIEs `aper:"optional"`
}
