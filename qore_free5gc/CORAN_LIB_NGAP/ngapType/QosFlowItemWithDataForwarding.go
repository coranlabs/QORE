package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type QosFlowItemWithDataForwarding struct {
	QosFlowIdentifier      QosFlowIdentifier
	DataForwardingAccepted *DataForwardingAccepted                                        `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerQosFlowItemWithDataForwardingExtIEs `aper:"optional"`
}
