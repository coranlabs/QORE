package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceToReleaseListHOCmd */
/* PDUSessionResourceToReleaseItemHOCmd */
type PDUSessionResourceToReleaseListHOCmd struct {
	List []PDUSessionResourceToReleaseItemHOCmd `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
