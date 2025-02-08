package factory

import (
	"fmt"
	"os"

	// "github.com/asaskevich/govalidator"
	"github.com/asaskevich/govalidator"
	"github.com/coranlabs/CORAN_SCP/Application_entity/logger"
	"gopkg.in/yaml.v2"
)

var ScpConfig *Config

func ReadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	if err := InitConfigFactory(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("ReadConfig [%s] Error: %+v", cfgPath, err)
	}
	if _, err := cfg.Validate(); err != nil {
		if validErrs, ok := err.(govalidator.Errors); ok {
			for _, validErr := range validErrs.Errors() {
				logger.CfgLog.Errorf("%+v", validErr)
			}
		} else {
			logger.CfgLog.Errorf("Validation error: %+v", err)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return nil, fmt.Errorf("Config validate Error")
	}
	

	return cfg, nil
}

func InitConfigFactory(f string, cfg *Config) error {
	if f == "" {
		// Use default config path
		f = ScpDefaultConfigPath
	}

	if content, err := os.ReadFile(f); err != nil {
		return fmt.Errorf("[Factory] %+v", err)
	} else {
		logger.CfgLog.Infof("Read config from [%s]", f)
		if yamlErr := yaml.Unmarshal(content, cfg); yamlErr != nil {
			return fmt.Errorf("[Factory] %+v", yamlErr)
		}
	}

	return nil
}
