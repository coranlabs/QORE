package message_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
	smf_pfcp "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp/message"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp/udp"
)

func TestSendPfcpAssociationSetupRequest(t *testing.T) {
}

func TestSendPfcpSessionEstablishmentResponse(t *testing.T) {
}

func TestSendPfcpSessionEstablishmentRequest(t *testing.T) {
}

func TestSendHeartbeatResponse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	smf_context.GetSelf().Ctx = ctx
	smf_context.GetSelf().PFCPCancelFunc = cancel
	udp.Run(smf_pfcp.Dispatch)

	udp.ServerStartTime = time.Now()
	var seq uint32 = 1
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 7001,
	}
	message.SendHeartbeatResponse(addr, seq)

	err := udp.ClosePfcp()
	require.NoError(t, err)
}
