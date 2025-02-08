package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type QosFlowAddOrModifyResponseItem struct {
	QosFlowIdentifier QosFlowIdentifier
	IEExtensions      *ProtocolExtensionContainerQosFlowAddOrModifyResponseItemExtIEs `aper:"optional"`
}
