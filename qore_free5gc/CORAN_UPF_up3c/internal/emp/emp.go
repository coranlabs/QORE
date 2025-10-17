package emp

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	// GTP-U message type for End Marker is 254
	EndMarkerMsgType = 254
	// GTP-U default port
	GTPUPort = 2152
)

// GTPUHeader represents the GTP-U Header
type GTPUHeader struct {
	Version     uint8
	PT          uint8
	MessageType uint8
	TEID        uint32
}

// Marshal converts the GTPUHeader struct to a byte slice
func (h *GTPUHeader) Marshal() []byte {
	header := make([]byte, 8)

	// First byte: Version (3 bits) + PT (1 bit) + Reserved (1 bit) + Extension Header flag (1 bit) + Sequence Number flag (1 bit) + N-PDU Number flag (1 bit)
	header[0] = (h.Version << 5) | (h.PT << 4)

	// Message type: End Marker (254)
	header[1] = h.MessageType

	// Length of the payload: 0 (since it's an End Marker, no payload)
	header[2] = 0
	header[3] = 0

	// TEID (Tunnel Endpoint Identifier) - 4 bytes
	header[4] = byte((h.TEID >> 24) & 0xFF)
	header[5] = byte((h.TEID >> 16) & 0xFF)
	header[6] = byte((h.TEID >> 8) & 0xFF)
	header[7] = byte(h.TEID & 0xFF)

	return header
}

func Send_emp(TEID uint32, gnbip net.IP) {
	// Get the destination IP from command-line arguments
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <destination-ip> <TEID>", os.Args[0])
	}
	// destIP := os.Args[1]
	// teid := os.Args[2]
	teid := TEID
	destIP := gnbip

	// Resolve the UDP address for GTP-U (UDP port 2152)
	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", destIP, GTPUPort))
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("Failed to connect to remote address: %v", err)
	}
	defer conn.Close()

	// Convert the TEID to uint32
	// var teidVal uint32
	// _, err = fmt.Sscanf(teid, "%x", &teidVal)
	// if err != nil {
	// 	log.Fatalf("Invalid TEID: %v", err)
	// }

	// Create a GTP-U End Marker packet
	header := GTPUHeader{
		Version:     1, // GTPv1-U
		PT:          1, // GTP-U
		MessageType: EndMarkerMsgType,
		TEID:        teid,
	}

	// Marshal the header into bytes
	gtpPacket := header.Marshal()

	// Send the GTP-U End Marker packet
	_, err = conn.Write(gtpPacket)
	if err != nil {
		log.Fatalf("Failed to send GTP-U End Marker packet: %v", err)
	}

	fmt.Printf("Sent GTP-U End Marker packet to %s with TEID %v\n", destIP, teid)
}
