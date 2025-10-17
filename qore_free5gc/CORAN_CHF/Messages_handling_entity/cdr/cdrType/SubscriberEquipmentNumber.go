package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SubscriberEquipmentNumber struct { /* Set Type */
	SubscriberEquipmentNumberType SubscriberEquipmentType `ber:"tagNum:0"`
	SubscriberEquipmentNumberData asn.OctetString         `ber:"tagNum:1"`
}
