package UPF_config

import (
	"log"

	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
)

var Conf UpfConfig

// Init init config for eupf package
func Initialize() {
	if err := Conf.Unmarshal(); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if err := Conf.Validate(); err != nil {
		log.Fatalf("eUPF config is invalid: %v", err)
	}

	logger.CfgLog.Debugf("Apply eUPF config: %+v", Conf)
}
