package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceHandoverList */
/* PDUSessionResourceHandoverItem */
type PDUSessionResourceHandoverList struct {
	List []PDUSessionResourceHandoverItem `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
