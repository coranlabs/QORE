package ngap_handler

import (
	"fmt"

	"github.com/coranlabs/CORAN_AMF/Application_entity/logger"
	"github.com/coranlabs/CORAN_AMF/Messages_controller/context"
	//
)

func SendToRan(ran *context.AmfRan, packet []byte) {
	defer func() {
		// This is workaround.
		// TODO: Handle ran.Conn close event correctly
		err := recover()
		if err != nil {
			logger.NgapLog.Warnf("Send error, gNB may have been lost: %+v", err)
		}
	}()

	if ran == nil {
		logger.NgapLog.Error("Ran is nil")
		return
	}

	if len(packet) == 0 {
		fmt.Printf("packet len is 0")
		return
	}

	if ran.Conn == nil {
		fmt.Printf("Ran conn is nil")
		return
	}

	if ran.Conn.RemoteAddr() == nil {
		fmt.Printf("Ran addr is nil")
		return
	}

	fmt.Printf("Send NGAP message To Ran")

	if n, err := ran.Conn.Write(packet); err != nil {
		fmt.Printf("Send error: %+v", err)
		return
	} else {
		fmt.Printf("Write %d bytes", n)
	}
}
