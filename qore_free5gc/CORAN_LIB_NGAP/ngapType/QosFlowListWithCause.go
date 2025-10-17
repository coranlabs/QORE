package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct QosFlowListWithCause */
/* QosFlowWithCauseItem */
type QosFlowListWithCause struct {
	List []QosFlowWithCauseItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
