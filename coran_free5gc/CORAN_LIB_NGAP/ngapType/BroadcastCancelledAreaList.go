package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

const (
	BroadcastCancelledAreaListPresentNothing int = iota /* No components present */
	BroadcastCancelledAreaListPresentCellIDCancelledEUTRA
	BroadcastCancelledAreaListPresentTAICancelledEUTRA
	BroadcastCancelledAreaListPresentEmergencyAreaIDCancelledEUTRA
	BroadcastCancelledAreaListPresentCellIDCancelledNR
	BroadcastCancelledAreaListPresentTAICancelledNR
	BroadcastCancelledAreaListPresentEmergencyAreaIDCancelledNR
	BroadcastCancelledAreaListPresentChoiceExtensions
)

type BroadcastCancelledAreaList struct {
	Present                       int
	CellIDCancelledEUTRA          *CellIDCancelledEUTRA
	TAICancelledEUTRA             *TAICancelledEUTRA
	EmergencyAreaIDCancelledEUTRA *EmergencyAreaIDCancelledEUTRA
	CellIDCancelledNR             *CellIDCancelledNR
	TAICancelledNR                *TAICancelledNR
	EmergencyAreaIDCancelledNR    *EmergencyAreaIDCancelledNR
	ChoiceExtensions              *ProtocolIESingleContainerBroadcastCancelledAreaListExtIEs
}
