package ngapType

// Need to import "github.com/coranlabs/CORAN_LIB_APER" if it uses "aper"

type UESecurityCapabilities struct {
	NRencryptionAlgorithms             NRencryptionAlgorithms
	NRintegrityProtectionAlgorithms    NRintegrityProtectionAlgorithms
	EUTRAencryptionAlgorithms          EUTRAencryptionAlgorithms
	EUTRAintegrityProtectionAlgorithms EUTRAintegrityProtectionAlgorithms
	IEExtensions                       *ProtocolExtensionContainerUESecurityCapabilitiesExtIEs `aper:"optional"`
}
