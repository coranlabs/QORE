package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ServiceAreaInformation */
/* ServiceAreaInformationItem */
type ServiceAreaInformation struct {
	List []ServiceAreaInformationItem `aper:"valueExt,sizeLB:1,sizeUB:16"`
}
