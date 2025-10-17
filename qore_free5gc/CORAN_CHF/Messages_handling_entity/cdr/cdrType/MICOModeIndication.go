package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	MICOModeIndicationPresentMICOMode   asn.Enumerated = 0
	MICOModeIndicationPresentNoMICOMode asn.Enumerated = 1
)

type MICOModeIndication struct {
	Value asn.Enumerated
}
