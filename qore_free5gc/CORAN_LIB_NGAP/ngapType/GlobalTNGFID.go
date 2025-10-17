package ngapType

// Need to import "gocoranlabs/lib/aper" if it uses "aper"

type GlobalTNGFID struct { /* Sequence Type (Extensible) */
	PLMNIdentity PLMNIdentity
	TNGFID       TNGFID                                        `aper:"valueLB:0,valueUB:1"`
	IEExtensions *ProtocolExtensionContainerGlobalTNGFIDExtIEs `aper:"optional"`
}
