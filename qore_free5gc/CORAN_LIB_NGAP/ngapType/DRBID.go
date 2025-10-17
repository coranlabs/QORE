package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type DRBID struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:32"`
}
