package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PriorityTypePresentLow    asn.Enumerated = 0
	PriorityTypePresentNormal asn.Enumerated = 1
	PriorityTypePresentHigh   asn.Enumerated = 2
)

type PriorityType struct {
	Value asn.Enumerated
}
