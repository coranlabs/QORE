package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	ThreeGPPPSDataOffStatusPresentActive   asn.Enumerated = 0
	ThreeGPPPSDataOffStatusPresentInactive asn.Enumerated = 1
)

type ThreeGPPPSDataOffStatus struct {
	Value asn.Enumerated
}
