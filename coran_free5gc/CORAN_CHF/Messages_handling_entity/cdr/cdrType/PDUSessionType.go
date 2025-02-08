package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	PDUSessionTypePresentIPv4v6       asn.Enumerated = 0
	PDUSessionTypePresentIPv4         asn.Enumerated = 1
	PDUSessionTypePresentIPv6         asn.Enumerated = 2
	PDUSessionTypePresentUnstructured asn.Enumerated = 3
	PDUSessionTypePresentEthernet     asn.Enumerated = 4
)

type PDUSessionType struct {
	Value asn.Enumerated
}
