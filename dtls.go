package dtls

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
	"log"
	"unsafe"
	// "net"
)

const (
	SSLMODE_SERVER = 0
	SSLMODE_CLIENT = 1
)

func check_error(err_code C.int) {
	if err_code != 1 {
		log.Fatalf("Error occurred: %d", C.ERR_get_error())
	}
}

func handle_ssl_error(ssl *C.SSL, err_code C.int) {

	switch err_code {
	case C.SSL_ERROR_WANT_READ:
		log.Fatal("SSL ERROR WANT READ")
	case C.SSL_ERROR_WANT_WRITE:
		log.Fatal("SSL ERROR WANT WRITE")
	case C.SSL_ERROR_ZERO_RETURN:
		log.Fatal("SSL connection closed\n")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SYSCALL:
		log.Fatal("SSL syscall error")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SSL:
		log.Fatal("SSL library error\n")
		C.SSL_free(ssl)
	default:
		log.Fatal("Unexpected SSL error\n")
		C.SSL_free(ssl)
	}
}

func Init_ssl_ctx(SSLMODE int, certPath string, keyPath string) *C.SSL_CTX {

	// C.SSL_library_init()
	// C.SSL_load_error_strings()

	// Create new SSL context using DTLS method
	ctx := C.SSL_CTX_new(C.DTLS_method())
	if ctx == nil {
		fmt.Println("Failed to create DTLS context")
		C.SSL_CTX_free(ctx)
		return nil
	}

	if SSLMODE == SSLMODE_SERVER {

		certPath := C.CString(certPath)
		keyPath := C.CString(keyPath)
		defer C.free(unsafe.Pointer(certPath))
		defer C.free(unsafe.Pointer(keyPath))

		if C.SSL_CTX_use_certificate_file(ctx, certPath, C.SSL_FILETYPE_PEM) <= 0 {
			log.Fatal("Failed to load certificate")
			C.SSL_CTX_free(ctx)
			return nil
		}

		if C.SSL_CTX_use_PrivateKey_file(ctx, keyPath, C.SSL_FILETYPE_PEM) <= 0 {
			log.Fatal("Failed to load private key")
			C.SSL_CTX_free(ctx)
			return nil
		}
	}

	return ctx

}

func New_ssl_conn(ctx *C.SSL_CTX, fd int, SSLMODE int) *C.SSL {

	// Create new DTLS connection using the context
	// var ssl *C.SSL
	ssl := C.SSL_new(ctx) // var ssl *C.SSL
	if ssl == nil {
		fmt.Println("Failed to create DTLS object")
		return nil
	}

	// Attach BIO to the SSL object
	// C.SSL_set_bio(ssl, bio, bio)
	C.SSL_set_fd(ssl, C.int(fd))

	if SSLMODE == SSLMODE_SERVER {
		C.SSL_set_accept_state(ssl)
	} else if SSLMODE == SSLMODE_CLIENT {
		C.SSL_set_connect_state(ssl)
	} else {
		log.Fatal("Wrong ssl mode passed. Pass either 0 or 1")
		return nil
	}

	// // Perform DTLS handshake
	// if res := C.SSL_connect(ssl); res != 1 {
	// 	check_error(res)
	// 	return
	// }
	return ssl
}

// call only when C.SSL_is_init_finished returns false
func Do_ssl_handshake(ssl *C.SSL) int {

	// if(!C.SSL_is_init_finished(ssl)){
	// }
	print_ssl_state(ssl)

	ret := C.SSL_do_handshake(ssl)

	if ret <= 0 {
		err_code := C.SSL_get_error(ssl, ret)
		print_ssl_state(ssl)
		handle_ssl_error(ssl, err_code)
		return -1
	}
	print_ssl_state(ssl)

	log.Println("SSL handshake done.")
	return 0

}

func print_ssl_state(ssl *C.SSL) {
	state := C.SSL_state_string_long(ssl)

	goState := C.GoString(state)

	// Print the state in Go
	log.Println("SSL state:", goState)

}

// decrypt_buf reads from the SSL connection into the destination buffer.
func decrypt_buf(ssl *C.SSL, dest []byte) C.int {
    // Convert Go byte slice to C pointer.
    cDest := (unsafe.Pointer(&dest[0]))
    n := C.SSL_read(ssl, cDest, C.int(len(dest)))
    return n
}

// func main() {

// 	// Create UDP connection (non-blocking UDP socket)
// 	conn, err := net.Dial("udp", "localhost:4444")
// 	if err != nil {
// 		fmt.Println("Failed to create UDP connection:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// // Convert Go socket to a BoringSSL BIO object for DTLS
// 	// bio := C.BIO_new_dgram(C.int(conn.(*net.UDPConn).Fd()), C.BIO_NOCLOSE)
// 	// if bio == nil {
// 	// 	fmt.Println("Failed to create DTLS BIO")
// 	// 	return
// 	// }

// 	file, err := conn.(*net.UDPConn).File()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fd := file.Fd()

// 	fmt.Println(fd)

// 	// bio := C.create_bio_from_fd(fd)
// 	// if bio == nil {
// 	// 	fmt.Println("Failed to create DTLS BIO")
// 	// 	return
// 	// }

// 	// fmt.Println("DTLS handshake successful!")
// }
