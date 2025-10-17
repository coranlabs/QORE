package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"

	circlPki "github.com/cloudflare/circl/pki"
	circlSign "github.com/cloudflare/circl/sign"
	"github.com/cloudflare/circl/sign/schemes"
)

var (
	ErrUnsupportedKeyType = errors.New("unsupported key type")
)

type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func ParsePrivateKeyFromPem(key []byte, circlType bool, schemeName string) (any, error) {

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	if !circlType {

		/* PKCS#1 outdated */
		if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			return key, nil
		}

		if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
			return key, nil
		}

		/* PKCS#8 is the general method */
		if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			switch key := key.(type) {
			case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
				return key, nil
			}
		}
	} else {

		// if schemeName != "" {

		// 	scheme := schemes.ByName(schemeName)

		// 	circlPriv, err := scheme.UnmarshalBinaryPrivateKey(block.Bytes)
		// 	if err == nil {
		// 		return circlPriv, nil
		// 	}

		// 	return nil, errors.New("failed to parse Circl private key")

		// } else {
		// 	return nil, errors.New("pass scheme name")
		// }

		var privKey pkcs8
		if _, err := asn1.Unmarshal(block.Bytes, &privKey); err != nil {
			return nil, err
		}

		scheme := circlPki.SchemeByOid(privKey.Algo.Algorithm)
		if scheme == nil {
			return nil, fmt.Errorf("x509: PKCS#8 wrapping contained private key with unknown algorithm: %v", privKey.Algo.Algorithm)
		}
		if l := len(privKey.Algo.Parameters.FullBytes); l != 0 {
			return nil, fmt.Errorf("x509: invalid %s private key parameters", scheme.Name())
		}
		var packedSk []byte
		if _, err := asn1.Unmarshal(privKey.PrivateKey, &packedSk); err != nil {
			return nil, fmt.Errorf("x509: invalid %s private key: %v", scheme.Name(), err)
		}
		sk, err := scheme.UnmarshalBinaryPrivateKey(packedSk)
		if err != nil {
			return nil, fmt.Errorf("x509: invalid %s private key: %v", scheme.Name(), err)
		}

		fmt.Println("parsed")

		return sk, nil
	}

	return nil, ErrUnsupportedKeyType

}

func ParsePublicKeyFromPem(key []byte, circlType bool, schemeName string) (any, error) {

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	if !circlType {

		/* unsure about PKIX for PQ keys yet */
		if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			switch key := key.(type) {
			case *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey, circlSign.PublicKey:
				return key, nil
			}
		}

		/*if the pub key is stored in an x509 certificate*/
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			return cert.PublicKey, nil
		}
	} else {

		if schemeName != "" {

			scheme := schemes.ByName(schemeName)

			circlPub, err := scheme.UnmarshalBinaryPublicKey(block.Bytes)
			if err == nil {
				return circlPub, nil
			}

			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				return cert.PublicKey, nil
			}

			return nil, errors.New("failed to parse Circl public key")

		} else {
			return nil, errors.New("pass scheme name")
		}
	}

	return nil, ErrUnsupportedKeyType

}
