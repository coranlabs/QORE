package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	RestrictionTypePresentAllowedAreas    asn.Enumerated = 0
	RestrictionTypePresentNotAllowedAreas asn.Enumerated = 1
)

type RestrictionType struct {
	Value asn.Enumerated
}
