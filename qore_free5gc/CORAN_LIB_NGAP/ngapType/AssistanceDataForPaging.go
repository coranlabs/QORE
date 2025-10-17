package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type AssistanceDataForPaging struct {
	AssistanceDataForRecommendedCells *AssistanceDataForRecommendedCells                       `aper:"valueExt,optional"`
	PagingAttemptInformation          *PagingAttemptInformation                                `aper:"valueExt,optional"`
	IEExtensions                      *ProtocolExtensionContainerAssistanceDataForPagingExtIEs `aper:"optional"`
}
