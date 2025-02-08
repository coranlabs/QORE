package data

import (
	"log"

	"github.com/coranlabs/HEXA_UPF/internal/logger"
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

	logger.MainLog.Tracef("Apply eUPF config: %+v", Conf)
}
