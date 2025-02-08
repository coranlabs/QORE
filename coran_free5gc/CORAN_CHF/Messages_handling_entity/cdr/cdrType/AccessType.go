package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	AccessTypePresentThreeGPPAccess    asn.Enumerated = 0
	AccessTypePresentNonThreeGPPAccess asn.Enumerated = 1
)

type AccessType struct {
	Value asn.Enumerated
}
