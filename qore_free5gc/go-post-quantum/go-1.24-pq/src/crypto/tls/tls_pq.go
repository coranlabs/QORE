package tls

import (
	"crypto/tls/internal/fips140tls"
	// "fmt"

	circlPki "github.com/cloudflare/circl/pki"
	circlSign "github.com/cloudflare/circl/sign"
	"github.com/cloudflare/circl/sign/eddilithium3"
	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
)

var PQSchemes = [...]struct {
	sigType uint8
	scheme  circlSign.Scheme
}{
	{signatureEdDilithium3, eddilithium3.Scheme()},
	{signatureMLDSA65, mldsa65.Scheme()},
}


func circlSchemeBySigType(sigType uint8) circlSign.Scheme {
	for _, cs := range PQSchemes {
		if cs.sigType == sigType {
			return cs.scheme
		}
	}
	return nil
}

func sigTypeByCirclScheme(scheme circlSign.Scheme) uint8 {
	for _, cs := range PQSchemes {
		if cs.scheme == scheme {
			return cs.sigType
		}
	}
	return 0
}

var supportedSignatureAlgorithmsWithPQ []SignatureScheme

// supportedSignatureAlgorithms returns enabled signature schemes. PQ signature
// schemes are only included when tls.Config#PQSignatureSchemesEnabled is set
// and FIPS-only mode is not enabled.
func (c *Config) supportedSignatureAlgorithms() []SignatureScheme {
	// If FIPS-only mode is requested, do not add other algos.
	if fips140tls.Required() {
		return supportedSignatureAlgorithms()
	}
	if !c.PQSignatureSchemesEnabled{
	}
	if c != nil && c.PQSignatureSchemesEnabled {
		return supportedSignatureAlgorithmsWithPQ
	}
	return defaultSupportedSignatureAlgorithms
}

func init() {
	supportedSignatureAlgorithmsWithPQ = append([]SignatureScheme{}, defaultSupportedSignatureAlgorithms...)
	for _, cs := range PQSchemes {
		supportedSignatureAlgorithmsWithPQ = append(supportedSignatureAlgorithmsWithPQ,
			SignatureScheme(cs.scheme.(circlPki.TLSScheme).TLSIdentifier()))
		// fmt.Printf("Added  signature scheme ID: 0x%04X\n", cs.sigType)
	}
}

