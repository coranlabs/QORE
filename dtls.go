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
	"os"
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
		fmt.Fprintf(os.Stderr, "Error occurred: %d", C.ERR_get_error())
	}
}

// a bit more layman
func handle_ssl_error(ssl *C.SSL, err_code C.int) C.int {

	switch err_code {
	case C.SSL_ERROR_WANT_READ:
		fmt.Fprintf(os.Stderr, "SSL ERROR WANT READ")
	case C.SSL_ERROR_WANT_WRITE:
		fmt.Fprintf(os.Stderr, "SSL ERROR WANT WRITE")
	case C.SSL_ERROR_ZERO_RETURN:
		fmt.Fprintf(os.Stderr, "SSL connection closed")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SYSCALL:
		fmt.Fprintf(os.Stderr, "SSL syscall error")
		C.SSL_free(ssl)
	case C.SSL_ERROR_SSL:
		fmt.Fprintf(os.Stderr, "SSL library error")
		C.SSL_free(ssl)
	default:
		fmt.Fprintf(os.Stderr, "Unexpected SSL error")
		C.SSL_free(ssl)
	}
	return err_code
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
			fmt.Fprintf(os.Stderr, "Failed to load certificate")
			C.SSL_CTX_free(ctx)
			return nil
		}

		if C.SSL_CTX_use_PrivateKey_file(ctx, keyPath, C.SSL_FILETYPE_PEM) <= 0 {
			fmt.Fprintf(os.Stderr, "Failed to load private key")
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
		fmt.Fprintf(os.Stderr, "Wrong ssl mode passed. Pass either 0 or 1")
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
		fmt.Fprintf(os.Stderr, "SSL BIO error.")
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

func queue_enc_bytes(ssl_conn *SSLConn, src []byte) int {

	srcLen := len(src)
	cSrc := unsafe.Pointer(&src[0])
	tbw := 0 //total bytes written

	// Write the encrypted data to the rbio
	for srcLen > 0 {
		bytesWritten := C.BIO_write(ssl_conn.rbio, cSrc, C.int(srcLen))
		if bytesWritten > 0 {
			tbw += int(bytesWritten)
			srcLen -= int(bytesWritten)
			cSrc = unsafe.Pointer(uintptr(cSrc) + uintptr(bytesWritten))
			log.Printf("Encrypted data written: %d bytes\n", bytesWritten)
		} else {
			log.Println("BIO_write: error or no data to write.")
			print_openssl_errors()
			return -1
		}
	}
	return tbw
}

func read_enc_buf(ssl_conn *SSLConn, src []byte, dest []byte) C.int {

	decryptedBytesLen := 0
	srcLenProc := 0 //src length processed

	// Perform SSL_read to decrypt data
	for {

		if srcLenProc < len(src) {
			bytesQueued := queue_enc_bytes(ssl_conn, src[srcLenProc:])
			if bytesQueued < 0 {
				fmt.Fprintf(os.Stderr, "Unable to queue the bytes")
				return -1 // Error during BIO_write
			}
			srcLenProc += bytesQueued
		} else {
			// return C.int(decryptedBytesLen)
			break
		}

		spaceRem := len(dest) - decryptedBytesLen
		if spaceRem <= 0 {
			log.Println("Destination buffer is full.")
			break
		}

		n := C.SSL_read(ssl_conn.ssl, unsafe.Pointer(&dest[decryptedBytesLen]), C.int(spaceRem))
		if n > 0 {
			decryptedBytesLen += int(n)
			log.Printf("Decrypted data: %d bytes\n", n)
		} else {
			errCode := C.SSL_get_error(ssl_conn.ssl, n)
			if errCode == C.SSL_ERROR_WANT_READ || errCode == C.SSL_ERROR_WANT_WRITE {
				log.Println("SSL_read requires more data or wants to write. Retrying...")
				continue
			} else {
				log.Printf("SSL_read error: %d\n", errCode)
				handle_ssl_error(ssl_conn.ssl, errCode)
				return -1
			}
		}
	}

	return C.int(decryptedBytesLen)
}

// a bit more technical
func print_openssl_errors() {
	for {
		err := C.ERR_get_error()
		if err == 0 {
			break
		}
		errStr := C.GoString(C.ERR_error_string(err, nil))
		log.Printf("OpenSSL Error: %s\n", errStr)
	}
}

func Encrypt_buf(ssl_conn *SSLConn, src []byte, dest []byte) C.int {
	srcLen := len(src)
	tbr := 0 // Total bytes read from wbio
	cSrc := unsafe.Pointer(&src[0])

	log.Printf("Encrypting data, src length: %d\n", srcLen)

	for srcLen > 0 {
		// Write plaintext data to SSL
		bytesWritten := C.SSL_write(ssl_conn.ssl, cSrc, C.int(srcLen))
		if bytesWritten <= 0 {
			errCode := C.SSL_get_error(ssl_conn.ssl, bytesWritten)
			if errCode == C.SSL_ERROR_WANT_READ || errCode == C.SSL_ERROR_WANT_WRITE {
				log.Println("SSL_write requires more data or wants to write. Retrying...")
				continue
			} else {
				log.Printf("SSL_write error: %d\n", errCode)
				handle_ssl_error(ssl_conn.ssl, errCode)
				return -1
			}
		}

		log.Printf("Bytes written to SSL: %d\n", int(bytesWritten))
		srcLen -= int(bytesWritten)
		log.Printf("srcLen: %d\n", srcLen)
		cSrc = unsafe.Pointer(uintptr(cSrc) + uintptr(bytesWritten))

		// Read encrypted data from wbio
		for {
			remSpace := len(dest) - tbr
			if remSpace <= 0 {
				log.Println("Destination buffer is full. Returning")
				return C.int(tbr)
			}

			bytesRead := C.BIO_read(ssl_conn.wbio, unsafe.Pointer(&dest[tbr]), C.int(remSpace))
			log.Println(bytesRead)
			if bytesRead > 0 {
				tbr += int(bytesRead)
				log.Printf("Encrypted data length: %d bytes", bytesRead)

			} else if bytesRead == 0 {
				log.Println("No more encrypted data in wbio.")
				break
			} else {
				log.Println("Error reading from BIO.")
				for {
					err := C.ERR_get_error()
					if err == 0 {
						log.Println("nil error")
						break
					}
					errStr := C.GoString(C.ERR_error_string(err, nil))
					log.Printf("OpenSSL Error: %s\n", errStr)
				}
				// print_openssl_errors()
				break
			}
		}
	}

	log.Printf("Total encrypted bytes: %d\n", tbr)
	return C.int(tbr)
}

func New_message_decrypt(ssl_conn *SSLConn, src []byte, dest []byte) C.int {

	n := read_enc_buf(ssl_conn, src, dest)
	if n < 0 {
		fmt.Fprintf(os.Stderr, "Error decrypting the message\n")
		return -1
	}
	return n

}
func New_message_encrypt(ssl_conn *SSLConn, src []byte, dest []byte) C.int {

	n := Encrypt_buf(ssl_conn, src, dest)
	if n < 0 {
		fmt.Fprintf(os.Stderr, "Error encrypting the message\n")
		return -1
	} else {
		return n
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
