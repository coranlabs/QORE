package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type DRBStatusDL12 struct {
	DLCOUNTValue COUNTValueForPDCPSN12                          `aper:"valueExt"`
	IEExtension  *ProtocolExtensionContainerDRBStatusDL12ExtIEs `aper:"optional"`
}
