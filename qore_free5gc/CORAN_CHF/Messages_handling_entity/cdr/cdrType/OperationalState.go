package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	OperationalStatePresentENABLED  asn.Enumerated = 0
	OperationalStatePresentDISABLED asn.Enumerated = 1
)

type OperationalState struct {
	Value asn.Enumerated
}
