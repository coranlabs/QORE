package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UEAggregateMaximumBitRate struct {
	UEAggregateMaximumBitRateDL BitRate
	UEAggregateMaximumBitRateUL BitRate
	IEExtensions                *ProtocolExtensionContainerUEAggregateMaximumBitRateExtIEs `aper:"optional"`
}
