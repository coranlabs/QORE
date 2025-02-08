// // // SPDX-License-Identifier: Apache-2.0
// // // Copyright 2024 CORAN LABS

package UPF_config

// import (
// 	"log"
// 	"time"

// 	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/sirupsen/logrus"
// 	"github.com/spf13/pflag"
// 	"github.com/spf13/viper"
// )

// // UpfConfig holds the configuration parameters for the UPF
// type UpfConfig struct {
// 	InterfaceName     []string      `mapstructure:"interface_name" json:"interface_name"`
// 	XDPAttachMode     string        `mapstructure:"xdp_attach_mode" validate:"oneof=generic native offload" json:"xdp_attach_mode"`
// 	ApiAddress        string        `mapstructure:"api_address" validate:"hostname_port" json:"api_address"`
// 	PfcpAddress       string        `mapstructure:"pfcp_address" validate:"hostname_port" json:"pfcp_address"`
// 	Gatewayip         string        `mapstructure:"Gatewayip" json:"Gatewayip"`
// 	PfcpNodeId        string        `mapstructure:"pfcp_node_id" validate:"hostname|ip" json:"pfcp_node_id"`
// 	MetricsAddress    string        `mapstructure:"metrics_address" validate:"hostname_port" json:"metrics_address"`
// 	N3Address         string        `mapstructure:"n3_address" validate:"ipv4" json:"n3_address"`
// 	GtpPeer           []string      `mapstructure:"gtp_peer" validate:"omitempty,dive,hostname_port" json:"gtp_peer"`
// 	EchoInterval      uint32        `mapstructure:"echo_interval" validate:"min=1" json:"echo_interval"`
// 	QerMapSize        uint32        `mapstructure:"qer_map_size" validate:"min=1" json:"qer_map_size"`
// 	FarMapSize        uint32        `mapstructure:"far_map_size" validate:"min=1" json:"far_map_size"`
// 	PdrMapSize        uint32        `mapstructure:"pdr_map_size" validate:"min=1" json:"pdr_map_size"`
// 	EbpfMapResize     bool          `mapstructure:"resize_ebpf_maps" json:"resize_ebpf_maps"`
// 	HeartbeatRetries  uint32        `mapstructure:"heartbeat_retries" json:"heartbeat_retries"`
// 	HeartbeatInterval uint32        `mapstructure:"heartbeat_interval" json:"heartbeat_interval"`
// 	HeartbeatTimeout  uint32        `mapstructure:"heartbeat_timeout" json:"heartbeat_timeout"`
// 	LoggingLevel      string        `mapstructure:"logging_level" validate:"required" json:"logging_level"`
// 	UEIPPool          string        `mapstructure:"ueip_pool" validate:"cidr" json:"ueip_pool"`
// 	FTEIDPool         uint32        `mapstructure:"teid_pool" json:"teid_pool"`
// 	FeatureUEIP       bool          `mapstructure:"feature_ueip" json:"feature_ueip"`
// 	FeatureFTUP       bool          `mapstructure:"feature_ftup" json:"feature_ftup"`
// 	MaxRetrans        uint32        `mapstructure:"max_retrans" validate:"min=1" json:"max_retrans"`         // New Field
// 	RetransTimeout    time.Duration `mapstructure:"retrans_timeout" validate:"min=1" json:"retrans_timeout"` // New Field
// }

// var (
// 	v    = viper.GetViper()
// 	Conf UpfConfig
// )

// func init() {
// 	setupLogger()
// 	registerFlags()
// 	bindFlags()
// 	setDefaults()
// 	readConfig()
// }

// func setupLogger() {
// 	logger.SetLogLevel(logrus.InfoLevel)
// }

// func registerFlags() {
// 	pflag.StringArray("iface", []string{}, "Interface list to bind XDP program to")
// 	pflag.String("attach", "generic", "XDP attach mode")
// 	pflag.String("aaddr", ":8080", "Address to bind API server to")
// 	pflag.String("paddr", "127.0.0.1:8805", "Address to bind PFCP server to")
// 	pflag.String("nodeid", "127.0.0.1", "PFCP Server Node ID")
// 	pflag.String("maddr", ":9090", "Address to bind metrics server to")
// 	pflag.String("n3addr", "127.0.0.1", "Address for communication over N3 interface")
// 	pflag.StringArray("peer", []string{}, "Address of GTP peer")
// 	pflag.Uint32("echo", 10, "Interval of sending echo requests in seconds")
// 	pflag.Uint32("qersize", 1024, "Size of the QER ebpf map")
// 	pflag.Uint32("farsize", 1024, "Size of the FAR ebpf map")
// 	pflag.Uint32("pdrsize", 1024, "Size of the PDR ebpf map")
// 	pflag.Bool("mapresize", false, "Enable or disable ebpf map resizing")
// 	pflag.Uint32("hbretries", 3, "Number of heartbeat retries")
// 	pflag.Uint32("hbinterval", 5, "Heartbeat interval in seconds")
// 	pflag.Uint32("hbtimeout", 5, "Heartbeat timeout in seconds")
// 	pflag.String("loglvl", "info", "Logging level")
// 	pflag.Bool("ueip", false, "Enable or disable UEIP feature")
// 	pflag.Bool("ftup", false, "Enable or disable FTUP feature")
// 	pflag.String("ueippool", "10.60.0.0/24", "IP pool for UEIP feature")
// 	pflag.Uint32("teidpool", 65535, "TEID pool for FTUP feature")
// 	pflag.String("Gatewayip", "10.100.50.236", "Gateway IP address")
// 	pflag.Uint32("maxretrans", 5, "Maximum number of retransmissions")   // New Flag
// 	pflag.Uint32("retranstimeout", 1000, "Retransmission timeout in ms") // New Flag
// 	pflag.Parse()
// }

// func bindFlags() {
// 	flags := map[string]string{
// 		"interface_name":     "iface",
// 		"xdp_attach_mode":    "attach",
// 		"api_address":        "aaddr",
// 		"pfcp_address":       "paddr",
// 		"pfcp_node_id":       "nodeid",
// 		"metrics_address":    "maddr",
// 		"n3_address":         "n3addr",
// 		"gtp_peer":           "peer",
// 		"echo_interval":      "echo",
// 		"qer_map_size":       "qersize",
// 		"far_map_size":       "farsize",
// 		"pdr_map_size":       "pdrsize",
// 		"resize_ebpf_maps":   "mapresize",
// 		"heartbeat_retries":  "hbretries",
// 		"heartbeat_interval": "hbinterval",
// 		"heartbeat_timeout":  "hbtimeout",
// 		"logging_level":      "loglvl",
// 		"feature_ueip":       "ueip",
// 		"feature_ftup":       "ftup",
// 		"ueip_pool":          "ueippool",
// 		"teid_pool":          "teidpool",
// 		"Gatewayip":          "Gatewayip",
// 	}

// 	for key, flag := range flags {
// 		_ = v.BindPFlag(key, pflag.Lookup(flag))
// 	}
// }

// func setDefaults() {
// 	defaults := map[string]interface{}{
// 		"interface_name":     "eth0",
// 		"xdp_attach_mode":    "generic",
// 		"api_address":        ":8080",
// 		"pfcp_address":       "10.100.200.14:8805",
// 		"pfcp_node_id":       "10.100.200.14",
// 		"metrics_address":    ":9090",
// 		"n3_address":         "10.100.200.14",
// 		"echo_interval":      10,
// 		"qer_map_size":       1024,
// 		"far_map_size":       1024,
// 		"pdr_map_size":       1024,
// 		"resize_ebpf_maps":   false,
// 		"heartbeat_retries":  3,
// 		"heartbeat_interval": 5,
// 		"heartbeat_timeout":  5,
// 		"logging_level":      "info",
// 		"feature_ueip":       false,
// 		"feature_ftup":       false,
// 		"ueip_pool":          "10.60.0.0/24",
// 		"teid_pool":          65535,
// 		"Gatewayip":          "10.100.50.236",
// 		"max_retrans":        5,    // Default value for MaxRetrans
// 		"retrans_timeout":    1000, // Default value for RetransTimeout (in milliseconds)
// 	}

// 	for key, value := range defaults {
// 		v.SetDefault(key, value)
// 	}
// }

// func readConfig() {
// 	v.SetConfigFile(pflag.Lookup("config").Value.String())
// 	v.SetEnvPrefix("upf")
// 	v.AutomaticEnv()

// 	if err := v.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			log.Print("Config file not found. Using defaults.")
// 		} else {
// 			logger.InitLog.Tracef("Unable to read config file: %v. Using defaults.", err)
// 		}
// 	}

// 	logger.InitLog.Tracef("Startup config: %+v", v.AllSettings())
// }

// func (c *UpfConfig) Validate() error {
// 	if err := validator.New().Struct(c); err != nil {
// 		return err
// 	}

// 	if !c.FeatureFTUP {
// 		c.FTEIDPool = 0
// 	}

// 	if !c.FeatureUEIP {
// 		c.UEIPPool = ""
// 	}

// 	return nil
// }

// func (c *UpfConfig) Unmarshal() error {
// 	return v.UnmarshalExact(c)
// }

// // Initialize initializes the UPF configuration
// func Initialize() {
// 	if err := Conf.Unmarshal(); err != nil {
// 		log.Fatalf("Unable to decode into struct: %v", err)
// 	}

// 	if err := Conf.Validate(); err != nil {
// 		log.Fatalf("eUPF config is invalid: %v", err)
// 	}

// 	logger.InitLog.Tracef("Apply eUPF config: %+v", Conf)
// }
