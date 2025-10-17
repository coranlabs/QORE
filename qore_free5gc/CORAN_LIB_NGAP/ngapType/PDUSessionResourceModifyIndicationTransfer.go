package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceModifyIndicationTransfer struct {
	DLQosFlowPerTNLInformation           QosFlowPerTNLInformation                                                    `aper:"valueExt"`
	AdditionalDLQosFlowPerTNLInformation *QosFlowPerTNLInformationList                                               `aper:"optional"`
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceModifyIndicationTransferExtIEs `aper:"optional"`
}
