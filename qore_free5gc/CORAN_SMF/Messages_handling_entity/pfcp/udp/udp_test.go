package udp_test

import (
	"context"
	"net"
	"testing"
	"time"

	pfcp "github.com/coranlabs/CORAN_LIB_PFCP"
	"github.com/stretchr/testify/require"

	"github.com/coranlabs/CORAN_LIB_PFCP/pfcpType"
	"github.com/coranlabs/CORAN_LIB_PFCP/pfcpUdp"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
	smf_pfcp "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp/udp"
)

const testPfcpClientPort = 12345

func TestRun(t *testing.T) {
	// Set SMF Node ID

	smf_context.GetSelf().CPNodeID = pfcpType.NodeID{
		NodeIdType: pfcpType.NodeIdTypeIpv4Address,
		IP:         net.ParseIP("127.0.0.1").To4(),
	}
	smf_context.GetSelf().ExternalAddr = "127.0.0.1"
	smf_context.GetSelf().ListenAddr = "127.0.0.1"

	ctx, cancel := context.WithCancel(context.Background())
	smf_context.GetSelf().Ctx = ctx
	smf_context.GetSelf().PFCPCancelFunc = cancel
	udp.Run(smf_pfcp.Dispatch)

	testPfcpReq := pfcp.Message{
		Header: pfcp.Header{
			Version:         1,
			MP:              0,
			S:               0,
			MessageType:     pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
			MessageLength:   9,
			SEID:            0,
			SequenceNumber:  1,
			MessagePriority: 0,
		},
		Body: pfcp.PFCPAssociationSetupRequest{
			NodeID: &pfcpType.NodeID{
				NodeIdType: 0,
				IP:         net.ParseIP("192.168.1.1").To4(),
			},
		},
	}

	srcAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: testPfcpClientPort,
	}
	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: pfcpUdp.PFCP_PORT,
	}

	err := pfcpUdp.SendPfcpMessage(testPfcpReq, srcAddr, dstAddr)
	require.Nil(t, err)

	err = udp.ClosePfcp()
	require.NoError(t, err)

	time.Sleep(300 * time.Millisecond)
}
