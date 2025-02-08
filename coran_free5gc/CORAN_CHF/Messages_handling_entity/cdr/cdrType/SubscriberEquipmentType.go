package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SubscriberEquipmentTypePresentIMEISV        asn.Enumerated = 0
	SubscriberEquipmentTypePresentMAC           asn.Enumerated = 1
	SubscriberEquipmentTypePresentEUI64         asn.Enumerated = 2
	SubscriberEquipmentTypePresentModifiedEUI64 asn.Enumerated = 3
)

type SubscriberEquipmentType struct {
	Value asn.Enumerated
}
