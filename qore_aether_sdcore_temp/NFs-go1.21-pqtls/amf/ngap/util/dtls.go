package util

import (
	"crypto/tls"
	"fmt"
	"net"
	// "time"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/pion/dtls/v3"
)

const (
	// DISABLED Client certificate authentication is disabled
	DISABLED string = "DISABLED"
	// WANTED Client certificate is requested but not required
	WANTED string = "WANTED"
	// NEEDED Client certificate is requested and required
	NEEDED string = "NEEDED"
)
const (
	// Server role
	Server string = "server"
	// Client role
	Client string = "client"
)

const (
	// Complete handshake and exit
	BASIC string = "BASIC"
	// Echo once and exit
	ONE_ECHO string = "ONE_ECHO"
	// Continuous echo
	FULL string = "FULL" //default
)

func loadCertificates(certFile string, keyFile string) *tls.Certificate {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	return &cert

}

func initDTLS(port int, serverName string, cipherSuiteID dtls.CipherSuiteID, cipherSuiteName string, certFile string, keyFile string) (*dtls.Config) {

	csMap := make(map[string]dtls.CipherSuiteID)
	csMap["TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"] = dtls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	csMap["TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"] = dtls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	csMap["TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"] = dtls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
	csMap["TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"] = dtls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA

	// create map of supported client authentication types
	caMap := make(map[string]dtls.ClientAuthType)
	caMap[DISABLED] = dtls.NoClientCert
	caMap[WANTED] = dtls.RequestClientCert
	caMap[NEEDED] = dtls.RequireAndVerifyClientCert


	var contains bool
	cipherSuiteID, contains = csMap[cipherSuiteName]
	if !contains {
		panic("Cipher suite " + cipherSuiteName + " not supported")
	}
	fmt.Println("Using cipher suite ", cipherSuiteName, " with id ", cipherSuiteID)

	var clientAuth dtls.ClientAuthType
	clientAuth = dtls.NoClientCert //default

	certificate := loadCertificates(certFile, keyFile)

	// // If a trusted certficate was provided, fetch it
	// var rootCAs *x509.CertPool = nil
	// if len(trustCert) > 0 {
	// 	dat, err := ioutil.ReadFile(trustCert)
	// 	util.Check(err)
	// 	rootCAs = x509.NewCertPool()
	// 	succ := rootCAs.AppendCertsFromPEM(dat)
	// 	if !succ {
	// 		panic("Was not successful in parsing certificate")
	// 	}
	// }

	// // If a client certificate was provided, fetch it, otherwise default to the trusted certificate
	// var clientCAs *x509.CertPool = nil
	// if len(clientCert) > 0 {
	// 	dat, err := ioutil.ReadFile(clientCert)
	// 	util.Check(err)
	// 	clientCAs = x509.NewCertPool()
	// 	succ := clientCAs.AppendCertsFromPEM(dat)
	// 	if !succ {
	// 		panic("Was not successful in parsing certificate")
	// 	}
	// } else {
	// 	clientCAs=rootCAs
	// }

	// certificates = make(certficate)
	var config *dtls.Config
	config = &dtls.Config{
		CipherSuites:         []dtls.CipherSuiteID{cipherSuiteID},
		ExtendedMasterSecret: dtls.DisableExtendedMasterSecret,
		Certificates:         []tls.Certificate{*certificate},
		ClientAuth:           clientAuth,
		ServerName:           serverName,
		InsecureHashes:       false,
		SignatureSchemes:     []tls.SignatureScheme{tls.PKCS1WithSHA1, tls.ECDSAWithSHA1, tls.ECDSAWithP256AndSHA256},
		//SignatureSchemes:     []tls.SignatureScheme{tls.PKCS1WithSHA1, tls.PKCS1WithSHA256, tls.PKCS1WithSHA384, tls.PKCS1WithSHA512, tls.ECDSAWithSHA1, tls.ECDSAWithP256AndSHA256, tls.ECDSAWithP384AndSHA384, tls.ECDSAWithP521AndSHA512},
	}

	return config

}

func startDTLSListener(config *dtls.Config,role string, sctpAddr *sctp.SCTPAddr){

	if role != Client && role != Server {
		panic("Role " + role + " is invalid")
	}

	udpAddr := net.ResolveUDPAddr()

	listener,err :=	dtls.Listen("udp",,config)




}