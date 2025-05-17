// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

//go:build !debug
// +build !debug

package http2_util

import (
	"crypto/tls"
	"fmt"
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
	default:
		return fmt.Sprintf("Unknown Curve ID: %d", id)
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
		Addr:    bindAddr,
		Handler: h2c.NewHandler(handler, h2Server),
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
					fmt.Println("No certs")
					return nil, err
				}
				fmt.Println("Cert found!")
				return &server.TLSConfig.Certificates[0], nil

			},
			CurvePreferences: []tls.CurveID{
				tls.X25519Kyber768Draft00, tls.X25519, tls.CurveP256},
			GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {

				fmt.Println(strings.Repeat("-", 30))
				fmt.Println("Client Details:\n")
				fmt.Printf("\tClient connected to: %s\n", chi.ServerName)

				fmt.Print("\tSupported Signature Schemes: ")
				for _, sigScheme := range chi.SignatureSchemes {
					fmt.Printf("\t\t%+v ", sigScheme)
				}
				fmt.Println()

				// Print the supported curves
				fmt.Println("\tSupported Curves: ")
				for _, curve := range chi.SupportedCurves {
					fmt.Printf("\t\t%+v ", curve)
				}
				fmt.Println()

				fmt.Println("\tCurve Preferences: ")
				for _, curve := range server.TLSConfig.CurvePreferences {
					fmt.Printf("\t\t%s\n", curveIDToString(curve))
				}
				fmt.Printf("\nServer Certificates: %d\n", len(server.TLSConfig.Certificates))

				return server.TLSConfig, nil

			},
		}
	}

	return
}
