package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const ( /* Enum Type */
	SMReplyPathRequestedPresentNoReplyPathSet asn.Enumerated = 0
	SMReplyPathRequestedPresentReplyPathSet   asn.Enumerated = 1
)

type SMReplyPathRequested struct {
	Value asn.Enumerated
}
