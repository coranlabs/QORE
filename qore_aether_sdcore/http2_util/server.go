// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

//go:build !debug
// +build !debug

package http2_util

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func curveIDToString(id tls.CurveID) string {
	switch id {
	case tls.CurveP256:
		return "P-256"
	case tls.CurveP384:
		return "P-384"
	case tls.CurveP521:
		return "P-521"
	case tls.X25519:
		return "X25519"
	case tls.X25519Kyber768Draft00:
		return "X25519-Kyber768-Draft00"
	case tls.P256Kyber768Draft00:
		return "P256-Kyber768-Draft00"
	case tls.MLKEM768:
		return "MLKEM768"
	case tls.SecP256r1MLKEM768:
		return "SecP256r1MLKEM768"
	case tls.X25519MLKEM768:
		return "X25519MLKEM768"
	default:
		return ""
	}
}

func printTLSHandshakeCipherSuite(conn *tls.Conn) {
	state := conn.ConnectionState()
	if state.HandshakeComplete {
		fmt.Printf("Handshake finished.\n")
		cipherSuite := tls.CipherSuiteName(state.CipherSuite)
		fmt.Printf("TLS version : %s\n", state.Version)
		fmt.Printf("Cipher Suite chosen : %s\n", cipherSuite)
		fmt.Printf("Negotiated protocol : %s\n", state.NegotiatedProtocol)
		fmt.Println()
	} else {

		err := conn.Handshake()
		if err != nil {
			log.Fatalf("Handshake not completed, Error: %s\n", err)
		} else {

			state := conn.ConnectionState()

			log.Println("Handshake done, Error: nil.")
			cipherSuite := tls.CipherSuiteName(state.CipherSuite)
			fmt.Printf("TLS version : %d\n", state.Version)
			fmt.Printf("Cipher Suite chosen : %s\n", cipherSuite)
			fmt.Printf("Negotiated protocol : %s\n", state.NegotiatedProtocol)
		}
	}
}

func serverConnStateHandler(conn net.Conn, state http.ConnState) {
	if state == http.StateNew {
		tlsConn, ok := conn.(*tls.Conn)
		if ok {
			fmt.Println("New PQ TLS connection established.")
			printTLSHandshakeCipherSuite(tlsConn)
		}
	} else if state == http.StateClosed {
		fmt.Println("Connection closed.")
	}
}

// NewServer returns a server instance with HTTP/2.0 and HTTP/2.0 cleartext support
// If this function cannot open or create the secret log file,
// **it still returns server instance** but without the secret log and error indication
func NewServer(bindAddr string, preMasterSecretLogPath string, handler http.Handler, cert tls.Certificate) (server *http.Server, err error) {
	if handler == nil {
		return nil, errors.New("server needs handler to handle request")
	}

	h2Server := &http2.Server{
		// TODO: extends the idle time after re-use openapi client
		IdleTimeout: 1 * time.Millisecond,
	}
	server = &http.Server{
		Addr:      bindAddr,
		Handler:   h2c.NewHandler(handler, h2Server),
		ConnState: serverConnStateHandler,
	}

	if preMasterSecretLogPath != "" {
		preMasterSecretFile, err := os.OpenFile(preMasterSecretLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return server, fmt.Errorf("create pre-master-secret log [%s] fail: %s", preMasterSecretLogPath, err)
		}
		server.TLSConfig = &tls.Config{
			KeyLogWriter:              preMasterSecretFile,
			PQSignatureSchemesEnabled: true,
			Certificates:              []tls.Certificate{cert},
			// PreferServerCipherSuites:  true, // deprecated - has no effect
			// MinVersion: tls.VersionTLS13,
			// ClientAuth: tls.NoClientCert,
			GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {

				if len(server.TLSConfig.Certificates) == 0 {
					fmt.Println("No certificates found for the SERVER.")
					return nil, err
				}
				fmt.Println("SERVER Certificate found.")
				fmt.Printf("\nNo. of Server Certificate(s): %d\n", len(server.TLSConfig.Certificates))

				return &server.TLSConfig.Certificates[0], nil

			},
			CurvePreferences: []tls.CurveID{tls.MLKEM768, tls.X25519MLKEM768, tls.SecP256r1MLKEM768,
				tls.X25519Kyber768Draft00, tls.X25519, tls.CurveP256},
			GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {

				fmt.Println(strings.Repeat("-", 30))
				fmt.Println("Connection Details:\n")

				fmt.Printf("\t1. Client Connected to: %s\n", chi.ServerName)

				fmt.Print("\t2. Client Supported Signature Schemes: ")
				for _, sigScheme := range chi.SignatureSchemes {
					fmt.Printf("\t\t%+v", sigScheme)
				}
				fmt.Println()

				// Print the supported curves
				fmt.Println("\t3. Client Supported Curves: ")
				for _, curve := range chi.SupportedCurves {
					curveString := curveIDToString(curve)
					if curveString == "" {
						fmt.Printf("\t\t•%+v\n", curve)
					} else {
						fmt.Printf("\t\t•%s\n", curveIDToString(curve))

					}
				}
				fmt.Println()

				fmt.Println("\t4. Curve Preferences: ")
				for _, curve := range server.TLSConfig.CurvePreferences {
					curveString := curveIDToString(curve)
					if curveString == "" {
						fmt.Printf("\t\t•%s\n", curve)
					} else {
						fmt.Printf("\t\t•%s\n", curveIDToString(curve))

					}
				}
				return server.TLSConfig, nil

			},
		}
	}

	return
}
