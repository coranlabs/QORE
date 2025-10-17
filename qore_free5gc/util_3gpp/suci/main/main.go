package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	// "hash"
	"log"
	"math"
)

func AnsiX963KDF(sharedKey, publicKey []byte, keyLenBytes int) []byte {

	// initializing a counter buffer of 4 bytes

	var counter uint32 = 0x00000001 //8 digits -> 4 bytes
	var kdfKey []byte

	// hash_len := sha256.Size
	outlen := 0

	for keyLenBytes > outlen {

		//make a byte slice of 4 bytes:
		counterBytes := make([]byte, 4)
		hasher := sha256.New()

		binary.BigEndian.PutUint32(counterBytes, counter)

		fmt.Printf("counterBytes: %x\n", counterBytes)

		// tmpK := sha256.Sum256(append(append(sharedKey,counterBytes...),publicKey...))
		// sliceK := tmpK[:]
		hasher.Write(sharedKey)
		hasher.Write(counterBytes)
		hasher.Write(publicKey)

		hash := hasher.Sum(nil)

		kdfKey = append(kdfKey, hash...)
		counter++
		outlen += len(hash)

	}
	fmt.Println("Size of KDF key: ", len(kdfKey))
	return kdfKey[0:keyLenBytes]

}

func AnsiX963KDF_2(sharedKey, publicKey []byte, profileEncKeyLen, profileMacKeyLen, profileHashLen int) []byte {
	var counter uint32 = 0x00000001
	var kdfKey []byte
	kdfRounds := int(math.Ceil(float64(profileEncKeyLen+profileMacKeyLen) / float64(profileHashLen)))
	for i := 1; i <= kdfRounds; i++ {
		counterBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(counterBytes, counter)
		fmt.Printf("counterBytes: %x\n", counterBytes)
		tmpK := sha256.Sum256(append(append(sharedKey, counterBytes...), publicKey...))
		sliceK := tmpK[:]
		kdfKey = append(kdfKey, sliceK...)
		// fmt.Printf("kdfKey in round %d: %x\n", i, kdfKey)
		counter++
	}
	return kdfKey
}

func HmacSha256(input, macKey []byte, macLen int) []byte {
	h := hmac.New(sha256.New, macKey)
	if _, err := h.Write(input); err != nil {
		log.Printf("HMAC SHA256 error %+v", err)
	}
	macVal := h.Sum(nil)
	macTag := macVal[:macLen]
	// fmt.Printf("macVal: %x\nmacTag: %x\n", macVal, macTag)
	return macTag
}

func main() {

	sharedKey := "BC498F7E3967A63CEBC881DE402D07E8D0A50E8F060ECE5B977733AAC8A93F6C"
	sharedKey_slice, _ := hex.DecodeString(sharedKey)

	publicKey := "005d9ae57b3269b8289fbc1989a8925642cae6a9063b77cafee412c962bc6e129c7335ad9c276c9f359641596c436ab541a4b14544a06f75aa191c4b621cc3234892b76106a3ccaa4e31a4336ab804c92fca10a5aa6b86318a1d14e47f4f848a6134c239a9824419aab14bbcc5968c47f367b7237402f178965c23b2970e259608aef95bdb73aeb862a09187912e457c42045b47c448a83902abf65784b84e1874ccc38a7f478aba0e03294b461e7e073d1ff0157579767faacb5c405479f31e310329bb90bf42956ab8bbb2f80a2473586a449080a2ac04d0a10be603a8e9841d31940d46593cae5a8357bb867acabf538bc3b94a6853794fdc573faa423014287ffe6a48efd296f628871cbc19adb8313c54b92b00bc5667ade839a057a090aa8087a70170e4ca39a0fc9e5516793e328c3a0a1ce2a10c87969262668a9073ccc4f102d9185ab307120d9c38eef5bf87712bfb553bc021573dd9c65ea06f7c1b783dacc0b283a4c56bc30438cc99c200fe13917de80414dc6e34496f8b472388e9b655948b11a354a5f043ff678c9ffcbcaf139f368cbb82d60c730c052df73ebd8a513f44c10641003a080eac8095f8ba8f4f662b6b70c7f8a8b6e8e752f37b7f27a28e7bcb128c196d70d05b4b31766cb2b23655c33d9598255c4e6f2963d3f15bb2887a1d471978e85f6babaa45f46152b4ca0a00183a8881dcb204ada26ba665b808329913963bd899c6b17c41c9f2801736686300a47fe02130f73d15f3c60b3675f90298a72a1ed4c1a3cd55108c42b39eb74c4a341ce005055ec212c7f380216a590c78b82a47685a77423046c5d00b70262ac041694a8db66893a18d06f9a4f362c08c55296c47ac6851514132ab7c293e44a3a6ae7c2fa4d36c776164024a9299245ff8684be7c5a4cb67104f8942bb24b453c07891e47ebeb2a69cf1c0e498cd0accb1a994894faacc61b5132a0c0332d965cdf231d32a9f3c85c5009091a54c6abe31abf3546f4554438da3a34421c60bc74a97e3914a6767057421e3f4a20982261321b7b3947287dc95e2e7860cac4bee3c22e1900826e82a0e8432e15ce70aa392bd0c5a7147a2ddb9588da38d31ccdb9c9ae8ee7c9bb9e2"

	pubKey_slice, _ := hex.DecodeString(publicKey)

	slice := AnsiX963KDF_2(sharedKey_slice, pubKey_slice, 32, 32, 32)

	fmt.Println("KDF 2: ")
	fmt.Printf("\nKey: %s\n", hex.EncodeToString(slice))

	macKey := slice[48:]

	fmt.Printf("\n ENC Key: %s\n", hex.EncodeToString(slice[0:32]))
	fmt.Printf("\n Mac key: %s\n", hex.EncodeToString(macKey))

	cipher_byte, _ := hex.DecodeString("E64D416B1F")

	fmt.Printf("Mac Tag: %s\n", hex.EncodeToString(HmacSha256(cipher_byte, macKey, 8)))

	fmt.Println(len(slice))
	fmt.Printf("\n\n")

	fmt.Println("KDF 1: ")
	slice_2 := AnsiX963KDF(sharedKey_slice, pubKey_slice, 80)

	fmt.Println()
	fmt.Printf("\nKey: %s\n", hex.EncodeToString(slice_2))

	macKey_2 := slice_2[48:]

	fmt.Printf("\n ENC Key: %s\n", hex.EncodeToString(slice_2[0:32]))
	fmt.Printf("\n Mac key: %s\n", hex.EncodeToString(macKey_2))

	cipher_byte_2, _ := hex.DecodeString("E64D416B1F")

	fmt.Printf("Mac Tag: %s\n", hex.EncodeToString(HmacSha256(cipher_byte_2, macKey_2, 8)))

}
