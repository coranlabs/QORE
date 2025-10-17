package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type QosFlowPerTNLInformationItem struct {
	QosFlowPerTNLInformation QosFlowPerTNLInformation                                      `aper:"valueExt"`
	IEExtensions             *ProtocolExtensionContainerQosFlowPerTNLInformationItemExtIEs `aper:"optional"`
}
