package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type AMFTNLAssociationToAddItem struct {
	AMFTNLAssociationAddress CPTransportLayerInformation `aper:"valueLB:0,valueUB:1"`
	TNLAssociationUsage      *TNLAssociationUsage        `aper:"optional"`
	TNLAddressWeightFactor   TNLAddressWeightFactor
	IEExtensions             *ProtocolExtensionContainerAMFTNLAssociationToAddItemExtIEs `aper:"optional"`
}
