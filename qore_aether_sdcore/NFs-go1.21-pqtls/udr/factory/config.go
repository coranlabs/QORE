// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

/*
 * UDR Configuration Factory
 */

package factory

import (
	protos "github.com/Nikhil690/connsert/proto/sdcoreConfig"
	"github.com/omec-project/logger_util"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/udr/logger"
)

const (
	UDR_EXPECTED_CONFIG_VERSION = "1.0.0"
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

const (
	UDR_DEFAULT_IPV4     = "127.0.0.4"
	UDR_DEFAULT_PORT     = "8000"
	UDR_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	Sbi             *Sbi              `yaml:"sbi"`
	Mongodb         *Mongodb          `yaml:"mongodb"`
	NrfUri          string            `yaml:"nrfUri"`
	PlmnSupportList []PlmnSupportItem `yaml:"plmnSupportList,omitempty"`
}

type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port"`
	Tls         *Tls   `yaml:"tls,omitempty"`
}

type Tls struct {
	Log string `yaml:"log"`
	Pem string `yaml:"pem"`
	Key string `yaml:"key"`
}

type Mongodb struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

var ConfigPodTrigger chan bool
var ConfigUpdateDbTrigger chan *UpdateDb

func init() {
	ConfigPodTrigger = make(chan bool)
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) addSmPolicyInfo(nwSlice *protos.NetworkSlice, dbUpdateChannel chan *UpdateDb) error {
	for _, devGrp := range nwSlice.DeviceGroup {
		for _, imsi := range devGrp.Imsi {
			smPolicyEntry := &SmPolicyUpdateEntry{
				Imsi:   imsi,
				Dnn:    devGrp.IpDomainDetails.DnnName,
				Snssai: nwSlice.Nssai,
			}
			dbUpdate := &UpdateDb{
				SmPolicyTable: smPolicyEntry,
			}
			dbUpdateChannel <- dbUpdate
		}
	}
	return nil
}

func (c *Config) updateConfig(commChannel chan *protos.NetworkSliceResponse, dbUpdateChannel chan *UpdateDb) bool {
	var minConfig bool
	for rsp := range commChannel {
		logger.GrpcLog.Infoln("Received updateConfig in the udr app : ")
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
					var found bool = false
					for _, cplmn := range UdrConfig.Configuration.PlmnSupportList {
						if (cplmn.PlmnId.Mnc == plmn.PlmnId.Mnc) && (cplmn.PlmnId.Mcc == plmn.PlmnId.Mcc) {
							found = true
							break
						}
					}
					if found == false {
						UdrConfig.Configuration.PlmnSupportList = append(UdrConfig.Configuration.PlmnSupportList, plmn)
					}
				} else {
					logger.GrpcLog.Infoln("Plmn not present in the message ")
				}

			}
			c.addSmPolicyInfo(ns, dbUpdateChannel)
		}
		if minConfig == false {
			// first slice Created
			if len(UdrConfig.Configuration.PlmnSupportList) > 0 {
				minConfig = true
				ConfigPodTrigger <- true
				logger.GrpcLog.Infoln("Send config trigger to main routine")
			}
		} else {
			// all slices deleted
			if len(UdrConfig.Configuration.PlmnSupportList) == 0 {
				minConfig = false
				ConfigPodTrigger <- false
				logger.GrpcLog.Infoln("Send config trigger to main routine")
			} else {
				ConfigPodTrigger <- true
				logger.GrpcLog.Infoln("Send config trigger to main routine")
			}
		}
	}
	return true
}
