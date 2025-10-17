package cdrType

import "github.com/coranlabs/CORAN_CHF/Messages_handling_entity/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	IPTextRepresentedAddressPresentNothing int = iota /* No components present */
	IPTextRepresentedAddressPresentIPTextV4Address
	IPTextRepresentedAddressPresentIPTextV6Address
)

type IPTextRepresentedAddress struct {
	Present         int            /* Choice Type */
	IPTextV4Address *asn.IA5String `ber:"tagNum:2"`
	IPTextV6Address *asn.IA5String `ber:"tagNum:3"`
}
