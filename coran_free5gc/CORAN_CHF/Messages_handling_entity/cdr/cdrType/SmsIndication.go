package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SmsIndicationPresentSMSSupported    asn.Enumerated = 0
	SmsIndicationPresentSMSNotSupported asn.Enumerated = 1
)

type SmsIndication struct {
	Value asn.Enumerated
}
