package utils

import (
	"context"
	"time"

	"github.com/coranlabs/CORAN_SMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/service"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp"
	"github.com/coranlabs/CORAN_SMF/Messages_handling_entity/pfcp/udp"
)

var (
	pfcpStart func(a *service.SmfApp)
	pfcpStop  func()
)

func InitPFCPFunc() (func(a *service.SmfApp), func()) {
	pfcpStart = func(a *service.SmfApp) {
		// Initialize PFCP server
		ctx, cancel := context.WithCancel(context.Background())
		smf_context.GetSelf().Ctx = ctx
		smf_context.GetSelf().PFCPCancelFunc = cancel

		udp.Run(pfcp.Dispatch)

		// Wait for PFCP start
		time.Sleep(1000 * time.Millisecond)

		for _, upNode := range smf_context.GetSelf().UserPlaneInformation.UPFs {
			upNode.UPF.Ctx, upNode.UPF.CancelFunc = context.WithCancel(ctx)
			go a.Processor().ToBeAssociatedWithUPF(ctx, upNode.UPF)
		}
	}

	pfcpStop = func() {
		smf_context.GetSelf().PFCPCancelFunc()
		err := udp.Server.Close()
		if err != nil {
			logger.Log.Errorf("udp server close failed %+v", err)
		}
	}

	return pfcpStart, pfcpStop
}
