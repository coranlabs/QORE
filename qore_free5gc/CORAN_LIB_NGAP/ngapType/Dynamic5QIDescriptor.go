package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type Dynamic5QIDescriptor struct {
	PriorityLevelQos       PriorityLevelQos
	PacketDelayBudget      PacketDelayBudget
	PacketErrorRate        PacketErrorRate                                       `aper:"valueExt"`
	FiveQI                 *FiveQI                                               `aper:"optional"`
	DelayCritical          *DelayCritical                                        `aper:"optional"`
	AveragingWindow        *AveragingWindow                                      `aper:"optional"`
	MaximumDataBurstVolume *MaximumDataBurstVolume                               `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerDynamic5QIDescriptorExtIEs `aper:"optional"`
}
