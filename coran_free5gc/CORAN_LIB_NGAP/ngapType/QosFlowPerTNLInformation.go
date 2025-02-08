package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type QosFlowPerTNLInformation struct {
	UPTransportLayerInformation UPTransportLayerInformation `aper:"valueLB:0,valueUB:1"`
	AssociatedQosFlowList       AssociatedQosFlowList
	IEExtensions                *ProtocolExtensionContainerQosFlowPerTNLInformationExtIEs `aper:"optional"`
}
