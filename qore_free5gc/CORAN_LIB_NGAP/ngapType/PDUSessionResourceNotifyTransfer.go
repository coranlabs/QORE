package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceNotifyTransfer struct {
	QosFlowNotifyList   *QosFlowNotifyList                                                `aper:"optional"`
	QosFlowReleasedList *QosFlowListWithCause                                             `aper:"optional"`
	IEExtensions        *ProtocolExtensionContainerPDUSessionResourceNotifyTransferExtIEs `aper:"optional"`
}
