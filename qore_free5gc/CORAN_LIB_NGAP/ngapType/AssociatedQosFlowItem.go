package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type AssociatedQosFlowItem struct {
	QosFlowIdentifier        QosFlowIdentifier
	QosFlowMappingIndication *aper.Enumerated                                       `aper:"optional,valueExt,valueLB:0,valueUB:1"`
	IEExtensions             *ProtocolExtensionContainerAssociatedQosFlowItemExtIEs `aper:"optional"`
}
