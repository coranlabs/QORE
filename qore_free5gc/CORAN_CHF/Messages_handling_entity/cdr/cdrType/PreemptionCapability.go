package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PreemptionCapabilityPresentNOTPREEMPT asn.Enumerated = 0
	PreemptionCapabilityPresentMAYPREEMPT asn.Enumerated = 1
)

type PreemptionCapability struct {
	Value asn.Enumerated
}
