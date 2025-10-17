package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type MaximumDataBurstVolume struct {
	Value int64 `aper:"valueExt,valueLB:0,valueUB:4095"`
}
