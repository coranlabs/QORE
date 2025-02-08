package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PDUSessionResourceFailedToSetupListCxtFail */
/* PDUSessionResourceFailedToSetupItemCxtFail */
type PDUSessionResourceFailedToSetupListCxtFail struct {
	List []PDUSessionResourceFailedToSetupItemCxtFail `aper:"valueExt,sizeLB:1,sizeUB:256"`
}
