package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <openssl/bio.h>
*/
import "C"

import (
	"fmt"
	"net"
)

func check_error(errCode C.int) {
	if errCode != 1 {
		fmt.Println("Error occurred:", C.ERR_get_error())
	}
}

func main() {
	// C.SSL_library_init()
	// C.SSL_load_error_strings()

	// Create new SSL context using DTLS method
	ctx := C.SSL_CTX_new(C.DTLS_method())
	if ctx == nil {
		fmt.Println("Failed to create DTLS context")
		return
	}
	defer C.SSL_CTX_free(ctx)

	// Create UDP connection (non-blocking UDP socket)
	conn, err := net.Dial("udp", "localhost:4444")
	if err != nil {
		fmt.Println("Failed to create UDP connection:", err)
		return
	}
	defer conn.Close()

	// // Convert Go socket to a BoringSSL BIO object for DTLS
	// bio := C.BIO_new_dgram(C.int(conn.(*net.UDPConn).Fd()), C.BIO_NOCLOSE)
	// if bio == nil {
	// 	fmt.Println("Failed to create DTLS BIO")
	// 	return
	// }

	file, err := conn.(*net.UDPConn).File()
	fd := file.Fd()

	// bio := C.create_bio_from_fd(fd)
	// if bio == nil {
	// 	fmt.Println("Failed to create DTLS BIO")
	// 	return
	// }

	// Create new DTLS connection using the context
	ssl := C.SSL_new(ctx)
	if ssl == nil {
		fmt.Println("Failed to create DTLS object")
		return
	}

	// Attach BIO to the SSL object
	// C.SSL_set_bio(ssl, bio, bio)
	C.SSL_set_fd(ssl, C.int(fd))

	// Perform DTLS handshake
	if res := C.SSL_connect(ssl); res != 1 {
		check_error(res)
		return
	}

	fmt.Println("DTLS handshake successful!")
	// Continue with reading/writing data..
}
