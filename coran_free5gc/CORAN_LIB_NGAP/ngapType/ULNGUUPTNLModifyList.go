package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UL_NGU_UP_TNLModifyList */
/* ULNGUUPTNLModifyItem */
type ULNGUUPTNLModifyList struct {
	List []ULNGUUPTNLModifyItem `aper:"valueExt,sizeLB:1,sizeUB:4"`
}
