/*
 * CHF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/asaskevich/govalidator"
	"github.com/coranlabs/CORAN_CHF/Application_entity/internal/logger"
)

var ChfConfig *Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string, cfg *Config) error {
	if f == "" {
		// Use default config path
		f = ChfDefaultConfigPath
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
func ReadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	if err := InitConfigFactory(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("ReadConfig [%s] Error: %+v", cfgPath, err)
	}
	if _, err := cfg.Validate(); err != nil {
		validErrsData, ok := err.(govalidator.Errors)
		if !ok {
			logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
			return nil, fmt.Errorf("Config validate Error: %+v", err)
		}
		for _, validErr := range validErrsData.Errors() {
			logger.CfgLog.Errorf("%+v", validErr)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return nil, fmt.Errorf("Config validate Error")
	}

	return cfg, nil
}
