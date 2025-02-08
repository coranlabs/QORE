package data

import (
	"log"

	"github.com/coranlabs/CORAN_UPF_eBPF/logger"
)

var Conf UpfConfig

// Init init config for eupf package
func Init() {
	if err := Conf.Unmarshal(); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if err := Conf.Validate(); err != nil {
		log.Fatalf("eUPF config is invalid: %v", err)
	}

	logger.InitLog.Tracef("Apply eUPF config: %+v", Conf)
}
