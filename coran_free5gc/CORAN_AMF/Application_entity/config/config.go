package config


type Amf_config struct {
	Name string `mapstructure:"name"`
	Addr string `mapstructure:"addr"`
	NumOstreams int
	MaxInstreams int
	MaxAttempts int
	MaxInitTimeout int
	
}