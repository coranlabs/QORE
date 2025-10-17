package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CompletedCellsInTAI_NR */
/* CompletedCellsInTAINRItem */
type CompletedCellsInTAINR struct {
	List []CompletedCellsInTAINRItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
