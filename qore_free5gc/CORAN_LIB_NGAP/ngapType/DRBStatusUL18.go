package ngapType

import aper "github.com/coranlabs/CORAN_LIB_APER"

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type DRBStatusUL18 struct {
	ULCOUNTValue              COUNTValueForPDCPSN18                          `aper:"valueExt"`
	ReceiveStatusOfULPDCPSDUs *aper.BitString                                `aper:"sizeLB:1,sizeUB:131072,optional"`
	IEExtension               *ProtocolExtensionContainerDRBStatusUL18ExtIEs `aper:"optional"`
}
