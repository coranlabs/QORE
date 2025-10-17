package jwt

import (
	"errors"

	"github.com/cloudflare/circl/sign"
	"github.com/cloudflare/circl/sign/schemes"
)

var (
	ErrCirclVerification = errors.New("circl scheme: verification error")
)

type SigningMethodCircl struct {
	scheme sign.Scheme
}

var SigningMethodPQ *SigningMethodCircl

func init() {
	defaultScheme := schemes.ByName("Ed448-Dilithium3")

	SigningMethodPQ = &SigningMethodCircl{scheme: defaultScheme}

	RegisterSigningMethod(SigningMethodPQ.Alg(), func() SigningMethod {
		return SigningMethodPQ
	})
}

func (m *SigningMethodCircl) Alg() string {
	return m.scheme.Name()
}

func (m *SigningMethodCircl) Sign(signingString string, key interface{}) ([]byte, error) {

	// var keyBytes []byte
	// switch k := key.(type) {
	// case []byte:
	// 	keyBytes = k
	// case string:
	// 	keyBytes = []byte(k)
	// case *dilithium.PrivateKey:
	// 	keyBytes = k.Bytes() // Extract private key bytes
	// default:
	// 	return nil, errors.New("invalid private key type")
	// }

	// privKey, ok := m.scheme.UnmarshalBinaryPrivateKey(keyBytes)
	// if !ok {
	// 	return nil, errors.New("invalid private key type")
	// }

	privKey, ok := key.(sign.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type")
	}

	sig := m.scheme.Sign(privKey, []byte(signingString), nil)
	return sig, nil
}

// Verify checks the signature using the selected scheme
func (m *SigningMethodCircl) Verify(signingString string, sig []byte, key interface{}) error {
	// pubKey, ok := m.scheme.UnmarshalBinaryPublicKey(key)
	// if !ok {
	// 	return errors.New("invalid public key type")
	// }

	pubKey, ok := key.(sign.PublicKey)

	if !ok {
		return errors.New("invalid private key type")
	}

	if !m.scheme.Verify(pubKey, []byte(signingString), sig, nil) {
		return ErrCirclVerification
	}
	return nil
}

// SetScheme allows dynamic selection of the signing scheme
func (m *SigningMethodCircl) SetScheme(schemeName string) error {
	scheme := schemes.ByName(schemeName)
	if scheme == nil {
		return errors.New("unsupported signing scheme")
	}
	m.scheme = scheme
	return nil
}
