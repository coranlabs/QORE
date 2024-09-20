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

type SSLConn struct {
	ssl         *C.SSL
	rbio        *C.BIO
	wbio        *C.BIO
	encrypt_buf []byte
}

const (
	SSLMODE_SERVER = 0
	SSLMODE_CLIENT = 1
)

func check_error(err_code C.int) {
	if err_code != 1 {
		fmt.Errorf("Error occurred: %d", C.ERR_get_error())
	}
}

func handle_ssl_error(ssl *C.SSL, err_code C.int) {

	switch err_code {
	case C.SSL_ERROR_WANT_READ:
		fmt.Errorf("SSL ERROR WANT READ")
	case C.SSL_ERROR_WANT_WRITE:
		fmt.Errorf("SSL ERROR WANT WRITE")
	case C.SSL_ERROR_ZERO_RETURN:
		fmt.Errorf("SSL connection closed")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SYSCALL:
		fmt.Errorf("SSL syscall error")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SSL:
		fmt.Errorf("SSL library error")
		C.SSL_free(ssl)
	default:
		fmt.Errorf("Unexpected SSL error")
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
			fmt.Errorf("Failed to load certificate")
			C.SSL_CTX_free(ctx)
			return nil
		}

		if C.SSL_CTX_use_PrivateKey_file(ctx, keyPath, C.SSL_FILETYPE_PEM) <= 0 {
			fmt.Errorf("Failed to load private key")
			C.SSL_CTX_free(ctx)
			return nil
		}
	}

	return ctx

}

func New_ssl_conn(ctx *C.SSL_CTX, fd int, SSLMODE int) *SSLConn {

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
		fmt.Errorf("Wrong ssl mode passed. Pass either 0 or 1")
		C.SSL_free(ssl)
		return nil
	}

	// // Perform DTLS handshake
	// if res := C.SSL_connect(ssl); res != 1 {
	// 	check_error(res)
	// 	return
	// }
	return &SSLConn{ssl: ssl}
}

// call only when C.SSL_is_init_finished returns false
func Do_ssl_handshake(ssl_conn *SSLConn) int {

	// if(!C.SSL_is_init_finished(ssl)){
	// }
	print_ssl_state(ssl_conn.ssl)

	ret := C.SSL_do_handshake(ssl_conn.ssl)

	if ret <= 0 {
		err_code := C.SSL_get_error(ssl_conn.ssl, ret)
		print_ssl_state(ssl_conn.ssl)
		handle_ssl_error(ssl_conn.ssl, err_code)
		return -1
	}
	print_ssl_state(ssl_conn.ssl)

	log.Println("SSL handshake done.")

	// set the rbio and wbio
	rbio := C.BIO_new(C.BIO_s_mem())
	// we write encrypted data from the socket to rbio using BIO write and then read from using read
	wbio := C.BIO_new(C.BIO_s_mem())

	if rbio == nil || wbio == nil {
		fmt.Errorf("SSL BIO error.")
		return -1
	}

	C.SSL_set_bio(ssl_conn.ssl, rbio, wbio)
	ssl_conn.rbio = rbio
	ssl_conn.wbio = wbio

	return 0

}

func print_ssl_state(ssl *C.SSL) {
	state := C.SSL_state_string_long(ssl)

	goState := C.GoString(state)

	// Print the state in Go
	log.Println("SSL state:", goState)

}

func read_enc_buf(ssl_conn *SSLConn, src []byte, dest []byte) C.int {
	src_len := len(src)
	cSrc := unsafe.Pointer(&src[0])
	cDest := unsafe.Pointer(&dest[0])

	for src_len > 0 {
		// Write the encrypted buffer to rbio
		bytes_written := C.BIO_write(ssl_conn.rbio, cSrc, C.int(src_len))
		if bytes_written <= 0 {
			return -1
		}

		log.Printf("Bytes written to the rbio: %d", bytes_written)
		src_len -= int(bytes_written)

		// Perform SSL read
		n := C.SSL_read(ssl_conn.ssl, cDest, C.int(len(dest)))
		if n <= 0 {
			return -1
		}

		// if src_len == 0 {
		// Convert C pointer (cDest) back to Go slice
		dec_bytes := C.GoBytes(cDest, n)

		copy(dest, dec_bytes[:n])

		log.Printf("Decrypted bytes copied: %d", n)
		return n
		// }

	}

	return 0
}

func Encrypt_buf(ssl_conn *SSLConn, src []byte, dest []byte) C.int {

	src_len := len(src)
	cSrc := (unsafe.Pointer(&src[0]))
	cDest := unsafe.Pointer(&dest[0])

	for src_len > 0 {
		n := C.SSL_write(ssl_conn.ssl, cSrc, C.int(src_len))
		log.Printf("Bytes written LENGTH: %d\n", int(n))

		if n <= 0 {
			fmt.Errorf("failed to encrypt the data\n")
			return -1
		}
		src_len -= int(n)
		bytes_read := C.BIO_read(ssl_conn.wbio, cDest, C.int(len(dest)))
		if bytes_read > 0 {
			log.Printf("Encrypted data length :%d\n", int(bytes_read))

			enc_bytes := C.GoBytes(cDest, bytes_read)

			copy(dest, enc_bytes[:bytes_read])

			return bytes_read
		}

	}
	return 0

}

func New_message_decrypt(ssl_conn *SSLConn, src []byte, dest []byte) {

	n := read_enc_buf(ssl_conn, src, dest)
	if n < 0 {
		fmt.Errorf("Error decrypting the message")
	}

}
func New_message_encrypt(ssl_conn *SSLConn, src []byte, dest []byte) {

	n := Encrypt_buf(ssl_conn, src, dest)
	if n < 0 {
		fmt.Errorf("Error encrypting the message")
	}

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
