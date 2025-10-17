package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type PDUSessionResourceInformationItem struct {
	PDUSessionID              PDUSessionID
	QosFlowInformationList    QosFlowInformationList
	DRBsToQosFlowsMappingList *DRBsToQosFlowsMappingList                                         `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceInformationItemExtIEs `aper:"optional"`
}
