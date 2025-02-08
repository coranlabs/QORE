package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UnavailableGUAMIList */
/* UnavailableGUAMIItem */
type UnavailableGUAMIList struct {
	List []UnavailableGUAMIItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
