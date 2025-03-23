// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

/*
 * NRF Configuration Factory
 */

package factory

import (
	"os"
	"strconv"

	protos "github.com/Nikhil690/connsert/proto/sdcoreConfig"
	"github.com/omec-project/logger_util"
	"github.com/omec-project/nrf/logger"
	"github.com/omec-project/openapi/models"
)

const (
	NRF_EXPECTED_CONFIG_VERSION = "1.0.0"
	NRF_DEFAULT_IPV4            = "127.0.0.10"
	NRF_DEFAULT_PORT            = "8000"
	NRF_DEFAULT_PORT_INT        = 8000
	NRF_DEFAULT_SCHEME          = "https"
	NRF_NFM_RES_URI_PREFIX      = "/nnrf-nfm/v1"
	NRF_DISC_RES_URI_PREFIX     = "/nnrf-disc/v1"
)

type Config struct {
	Info          *Info               `yaml:"info"`
	Configuration *Configuration      `yaml:"configuration"`
	Logger        *logger_util.Logger `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Configuration struct {
	Sbi                   *Sbi              `yaml:"sbi,omitempty"`
	MongoDBName           string            `yaml:"MongoDBName"`
	MongoDBUrl            string            `yaml:"MongoDBUrl"`
	MongoDBStreamEnable   bool              `yaml:"mongoDBStreamEnable"`
	NfProfileExpiryEnable bool              `yaml:"nfProfileExpiryEnable"`
	DefaultPlmnId         models.PlmnId     `yaml:"DefaultPlmnId"`
	ServiceNameList       []string          `yaml:"serviceNameList,omitempty"`
	PlmnSupportList       []PlmnSupportItem `yaml:"plmnSupportList,omitempty"`
	NfKeepAliveTime       int32             `yaml:"nfKeepAliveTime,omitempty"`
}

type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is serviced or registered at another NRF.
	// IPv6Addr  string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
}

var MinConfigAvailable bool

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) GetSbiScheme() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Scheme != "" {
		return c.Configuration.Sbi.Scheme
	}
	return NRF_DEFAULT_SCHEME
}

func (c *Config) GetSbiPort() int {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Port != 0 {
		return c.Configuration.Sbi.Port
	}
	return NRF_DEFAULT_PORT_INT
}

func (c *Config) GetSbiBindingAddr() string {
	var bindAddr string
	if c.Configuration == nil || c.Configuration.Sbi == nil {
		return "0.0.0.0:" + NRF_DEFAULT_PORT
	}
	if c.Configuration.Sbi.BindingIPv4 != "" {
		if bindIPv4 := os.Getenv(c.Configuration.Sbi.BindingIPv4); bindIPv4 != "" {
			logger.CfgLog.Infof("Parsing ServerIPv4 [%s] from ENV Variable", bindIPv4)
			bindAddr = bindIPv4 + ":"
		} else {
			bindAddr = c.Configuration.Sbi.BindingIPv4 + ":"
		}
	} else {
		bindAddr = "0.0.0.0:"
	}
	if c.Configuration.Sbi.Port != 0 {
		bindAddr = bindAddr + strconv.Itoa(c.Configuration.Sbi.Port)
	} else {
		bindAddr = bindAddr + NRF_DEFAULT_PORT
	}
	return bindAddr
}

func (c *Config) GetSbiRegisterIP() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.RegisterIPv4 != "" {
		return c.Configuration.Sbi.RegisterIPv4
	}
	return NRF_DEFAULT_IPV4
}

func (c *Config) GetSbiRegisterAddr() string {
	regAddr := c.GetSbiRegisterIP() + ":"
	if c.Configuration.Sbi.Port != 0 {
		regAddr = regAddr + strconv.Itoa(c.Configuration.Sbi.Port)
	} else {
		regAddr = regAddr + NRF_DEFAULT_PORT
	}
	return regAddr
}

func (c *Config) GetSbiUri() string {
	return c.GetSbiScheme() + "://" + c.GetSbiRegisterAddr()
}

func (c *Config) updateConfig(commChannel chan *protos.NetworkSliceResponse) bool {
	for rsp := range commChannel {
		logger.GrpcLog.Infoln("Received updateConfig in the nrf app : ")
		logger.GrpcLog.Info("+---------------------------------------------+")
		logger.GrpcLog.Infof("| %-43s |\n", "Network Slice")
		logger.GrpcLog.Infof("|---------------------------------------------|")
		// logger.GrpcLog.Infof("| %15s | %10d |\n", "RestartCounter", rsp.RestartCounter)
		// logger.GrpcLog.Infof("| %15s | %10d |\n", "ConfigUpdated", rsp.ConfigUpdated)
		for _, slice := range rsp.NetworkSlice {
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Name", slice.Name)
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Sst", slice.Nssai.Sst)
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Sd", slice.Nssai.Sd)
			logger.GrpcLog.Infof("|---------------------------------------------|")
			for _, group := range slice.DeviceGroup {
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Device Group", group.Name)
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "IP Domain Details", group.IpDomainDetails.Name)
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "DNN Name", group.IpDomainDetails.DnnName)
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "UE Pool", group.IpDomainDetails.UePool)
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "DNS Primary", group.IpDomainDetails.DnsPrimary)
				logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "MTU", group.IpDomainDetails.Mtu)
				logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "DnnMbrUplink", group.IpDomainDetails.UeDnnQos.DnnMbrUplink)
				logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "DnnMbrDownlink", group.IpDomainDetails.UeDnnQos.DnnMbrDownlink)
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Traffic Class", group.IpDomainDetails.UeDnnQos.TrafficClass.Name)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "QCI", group.IpDomainDetails.UeDnnQos.TrafficClass.Qci)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "ARP", group.IpDomainDetails.UeDnnQos.TrafficClass.Arp)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "PDB", group.IpDomainDetails.UeDnnQos.TrafficClass.Pdb)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "PELR", group.IpDomainDetails.UeDnnQos.TrafficClass.Pelr)
				// for _, imdetails := range group.Imsi {
				// 	logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "IMSI Supported", imdetails)
				// }

				for i, imdetails := range group.Imsi {
					label := ""
					if i == len(group.Imsi)/2 {
						label = "IMSI_Supported"
					}
					logger.GrpcLog.Infof("| %-18s  | %-21s |\n", label, imdetails)
				}
				logger.GrpcLog.Info("|---------------------------------------------|")
			}
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Site", slice.Site.SiteName)
			for _, gnb := range slice.Site.Gnb {
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "GNB", gnb.Name)
				logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "TAC", gnb.Tac)
				logger.GrpcLog.Info("|---------------------------------------------|")
			}
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "MCC", slice.Site.Plmn.Mcc)
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "MNC", slice.Site.Plmn.Mnc)
			logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "UPF", slice.Site.Upf.UpfName)
			for _, appfilter := range slice.AppFilters.PccRuleBase {
				for _, flowinfo := range appfilter.FlowInfos {
					// logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Flow Description", flowinfo.FlowDesc)
					logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Traffic Class", flowinfo.TosTrafficClass)
					logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Flow Direction", flowinfo.FlowDir)
					logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Flow Status", flowinfo.FlowStatus)
				}
				logger.GrpcLog.Infof("| %-18s  | %-21s |\n", "Rule ID", appfilter.RuleId)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "Var5qi", appfilter.Qos.Var5Qi)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "ARP:PL", appfilter.Qos.Arp.PL)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "ARP:PC", appfilter.Qos.Arp.PC)
				// logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "ARP:PV", appfilter.Qos.Arp.PV)
				logger.GrpcLog.Infof("| %-18s  | %-21d |\n", "Priority", appfilter.Priority)
			}
			logger.GrpcLog.Info("+---------------------------------------------+")
		}
		for _, ns := range rsp.NetworkSlice {
			logger.GrpcLog.Infoln("Network Slice Name ", ns.Name)
			if ns.Site != nil {
				logger.GrpcLog.Infoln("Network Slice has site name present ")
				site := ns.Site
				logger.GrpcLog.Infoln("Site name ", site.SiteName)
				if site.Plmn != nil {
					logger.GrpcLog.Infoln("Plmn mcc ", site.Plmn.Mcc)
					plmn := PlmnSupportItem{}
					plmn.PlmnId.Mnc = site.Plmn.Mnc
					plmn.PlmnId.Mcc = site.Plmn.Mcc
					NrfConfig.Configuration.PlmnSupportList = append(NrfConfig.Configuration.PlmnSupportList, plmn)
				} else {
					logger.GrpcLog.Infoln("Plmn not present in the message ")
				}

			}
		}
		logger.GrpcLog.Infoln("minimum config Available")
		MinConfigAvailable = true
	}
	return true
}
