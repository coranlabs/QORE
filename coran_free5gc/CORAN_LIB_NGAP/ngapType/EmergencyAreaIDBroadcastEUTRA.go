package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct EmergencyAreaIDBroadcastEUTRA */
/* EmergencyAreaIDBroadcastEUTRAItem */
type EmergencyAreaIDBroadcastEUTRA struct {
	List []EmergencyAreaIDBroadcastEUTRAItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
