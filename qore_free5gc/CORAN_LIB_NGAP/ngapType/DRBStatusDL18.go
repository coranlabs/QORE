package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type DRBStatusDL18 struct {
	DLCOUNTValue COUNTValueForPDCPSN18                          `aper:"valueExt"`
	IEExtension  *ProtocolExtensionContainerDRBStatusDL18ExtIEs `aper:"optional"`
}
