package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct SliceOverloadList */
/* SliceOverloadItem */
type SliceOverloadList struct {
	List []SliceOverloadItem `aper:"valueExt,sizeLB:1,sizeUB:1024"`
}
