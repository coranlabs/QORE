// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
/*
 * AMF Unit Testcases
 *
 */
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	protos "github.com/Nikhil690/connsert/proto/sdcoreConfig"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/pcf/consumer"
	"github.com/omec-project/pcf/context"
	"github.com/omec-project/pcf/service"
	"github.com/stretchr/testify/require"
)

var PCFTest = &service.PCF{}
var bitRateValues = make(map[int64]string)

func init() {
	bitRateValues = map[int64]string{
		1000:        "1 Kbps",
		67200:       "67 Kbps",
		777111:      "777 Kbps",
		77711000:    "77 Mbps",
		64435000:    "64435 Kbps",
		77711000000: "77 Gbps",
		64435000111: "64435 Mbps",
	}
}

/*
	func init() {
		factory.InitConfigFactory("amfTest/amfcfg.yaml")
	}
*/
func GetNetworkSliceConfig() *protos.NetworkSliceResponse {
	var rsp protos.NetworkSliceResponse
	rsp.NetworkSlice = make([]*protos.NetworkSlice, 0)
	ns := protos.NetworkSlice{}
	slice := protos.NSSAI{Sst: "1", Sd: "010203"}
	ns.Nssai = &slice
	site := protos.SiteInfo{SiteName: "siteOne", Gnb: make([]*protos.GNodeB, 0), Plmn: new(protos.PlmnId)}
	gNb := protos.GNodeB{Name: "gnb", Tac: 1}
	site.Gnb = append(site.Gnb, &gNb)
	site.Plmn.Mcc = "208"
	site.Plmn.Mnc = "93"
	ns.Site = &site
	rsp.NetworkSlice = append(rsp.NetworkSlice, &ns)
	return &rsp
}

/*func TestInitialConfig(t *testing.T) {
	factory.AmfConfig.Configuration.PlmnSupportList = nil
	factory.AmfConfig.Configuration.ServedGumaiList = nil
	factory.AmfConfig.Configuration.SupportTAIList = nil
	var Rsp chan *protos.NetworkSliceResponse
	Rsp = make(chan *protos.NetworkSliceResponse)
	go func() {
		Rsp <- GetNetworkSliceConfig()
	}()
	go func() {
		AMF.UpdateConfig(Rsp)
	}()
	time.Sleep(2 * time.Second)
	if factory.AmfConfig.Configuration.PlmnSupportList != nil &&
		factory.AmfConfig.Configuration.ServedGumaiList != nil &&
		factory.AmfConfig.Configuration.SupportTAIList != nil {
		fmt.Printf("test passed")
	} else {
		t.Errorf("test failed")
	}
}*/
// data in JSON format which
// is to be decoded
var Data = []byte(`{
	"NetworkSlice": [
		{
		 "Name": "siteOne",
		 "Nssai": {"Sst": "010203", "Sd": "1"},
		 "Site": {
			"SiteName": "siteOne",
			"Gnb": [
				{"Name": "gnb1", "Tac": 1}, 
				{"Name": "gnb2", "Tac": 2}
			],
			"Plmn": {"mcc": "208", "mnc": "93"}
		  },
		 "DeviceGroup": [
		 	    {"Name": "dg1", 
				"IpDomainDetails": {
				"DnnName": "internet",
				"UeDnnQos": {"DnnMbrUplink": 100000, "DnnMbrDownlink": 50000, "TrafficClass": {"Qci": 9, "Arp": 6}}
				},
				"Imsi": ["123456789123456", "123456789123457", "123456789123458"]}
		 ],
		 "AppFilters": {
			"PccRuleBase": [{"FlowInfos": [{"FlowDesc": "permit out ip from 8.8.8.8/32 to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule1", "Priority": 15, "Qos": {"Var5qi": 8}},
						   {"FlowInfos": [{"FlowDesc": "permit out ip from any to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule2", "Priority": 25}]
		 },
		 "OperationType": 0
		}
		]}`)

/*
	func TestUpdateConfig(t *testing.T) {
		var nrp protos.NetworkSliceResponse
		err := json.Unmarshal(Data, &nrp)
		if err != nil {
			panic(err)
		}
		var Rsp chan *protos.NetworkSliceResponse
		Rsp = make(chan *protos.NetworkSliceResponse)
		go func() {
			Rsp <- &nrp
		}()
		go func() {
			AMF.UpdateConfig(Rsp)
		}()
		time.Sleep(2 * time.Second)
		if factory.AmfConfig.Configuration.SupportTAIList != nil &&
			len(factory.AmfConfig.Configuration.SupportTAIList) == 2 {
			fmt.Printf("test passed")
		} else {
			t.Errorf("test failed")
		}
	}
*/
func TestUpdatePcfSubsriberPolicyDataAdd(t *testing.T) {
	var nrp protos.NetworkSliceResponse
	err := json.Unmarshal(Data, &nrp)
	if err != nil {
		panic(err)
	}
	for _, ns := range nrp.NetworkSlice {
		PCFTest.UpdatePcfSubsriberPolicyData(ns)
	}
	self := context.PCF_Self()
	if len(self.PcfSubscriberPolicyData) == 3 {
		fmt.Printf("test case TestUpdatePcfSubsriberPolicyDataAdd Passed\n")
	} else {
		t.Errorf("test case failed\n")
	}
}

var UData = []byte(`{
	"NetworkSlice": [
		{
		 "Name": "siteOne",
		 "Nssai": {"Sst": "010203", "Sd": "1"},
		 "Site": {
			"SiteName": "siteOne",
			"Gnb": [
				{"Name": "gnb1", "Tac": 1}, 
				{"Name": "gnb2", "Tac": 2}
			],
			"Plmn": {"mcc": "208", "mnc": "93"}
		  },
		 "DeviceGroup": [
		 	    {"Name": "dg1", 
				"IpDomainDetails": {
				"DnnName": "internet",
				"UeDnnQos": {"DnnMbrUplink": 100000, "DnnMbrDownlink": 50000, "TrafficClass": {"Qci": 9, "Arp": 6}}
				},
				"Imsi": ["123456789123456", "123456789123457", "123456789123458", "123456789123459", "123456789123460"]}
		 ],
		 "AppFilters": {
			"PccRuleBase": [{"FlowInfos": [{"FlowDesc": "permit out ip from 8.8.8.8/32 to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule1", "Priority": 15, "Qos": {"Var5qi": 8}},
						   {"FlowInfos": [{"FlowDesc": "permit out ip from any to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule2", "Priority": 25}]
		 },
		 "OperationType": 1,
		 "AddUpdatedImsis": ["123456789123459", "123456789123460"]
		}
		]}`)

func TestUpdatePcfSubsriberPolicyDataUpdate(t *testing.T) {
	var nrp protos.NetworkSliceResponse
	err := json.Unmarshal(UData, &nrp)
	if err != nil {
		panic(err)
	}
	for _, ns := range nrp.NetworkSlice {
		PCFTest.UpdatePcfSubsriberPolicyData(ns)
	}
	self := context.PCF_Self()
	if len(self.PcfSubscriberPolicyData) == 5 {
		fmt.Printf("test case TestUpdatePcfSubsriberPolicyDataUpdate Passed\n")
	} else {
		t.Errorf("test case failed\n")
	}
}

var UData1 = []byte(`{
	"NetworkSlice": [
		{
		 "Name": "siteOne",
		 "Nssai": {"Sst": "010203", "Sd": "1"},
		 "Site": {
			"SiteName": "siteOne",
			"Gnb": [
				{"Name": "gnb1", "Tac": 1}, 
				{"Name": "gnb2", "Tac": 2}
			],
			"Plmn": {"mcc": "208", "mnc": "93"}
		  },
		 "DeviceGroup": [
		 	    {"Name": "dg1", 
				"IpDomainDetails": {
				"DnnName": "internet",
				"UeDnnQos": {"DnnMbrUplink": 100000, "DnnMbrDownlink": 50000, "TrafficClass": {"Qci": 9, "Arp": 6}}
				},
				"Imsi": ["123456789123456", "123456789123459", "123456789123460", "123456789123461"]}
		 ],
		 "AppFilters": {
			"PccRuleBase": [{"FlowInfos": [{"FlowDesc": "permit out ip from 8.8.8.8/32 to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule1", "Priority": 15, "Qos": {"Var5qi": 8}},
						   {"FlowInfos": [{"FlowDesc": "permit out ip from any to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule2", "Priority": 25}]
		 },
		 "OperationType": 1,
		 "AddUpdatedImsis": ["123456789123461"],
		 "DeletedImsis": ["123456789123457", "123456789123458"]
		}
		]}`)

// Two imsis deleted and 1 imsi added in device group
func TestUpdatePcfSubsriberPolicyDataUpdate1(t *testing.T) {
	var nrp protos.NetworkSliceResponse
	err := json.Unmarshal(UData1, &nrp)
	if err != nil {
		panic(err)
	}
	for _, ns := range nrp.NetworkSlice {
		PCFTest.UpdatePcfSubsriberPolicyData(ns)
	}
	self := context.PCF_Self()
	if len(self.PcfSubscriberPolicyData) == 4 {
		fmt.Printf("test case TestUpdatePcfSubsriberPolicyDataUpdate1 Passed\n")
	} else {
		t.Errorf("test case failed\n")
	}
}

var DelData = []byte(`{
	"NetworkSlice": [
		{
		 "Name": "siteOne",
		 "Nssai": {"Sst": "010203", "Sd": "1"},
		 "Site": {
			"SiteName": "siteOne",
			"Gnb": [
				{"Name": "gnb1", "Tac": 1}, 
				{"Name": "gnb2", "Tac": 2}
			],
			"Plmn": {"mcc": "208", "mnc": "93"}
		  },
		 "DeviceGroup": [
		 	    {"Name": "dg1", 
				"IpDomainDetails": {
				"DnnName": "internet",
				"UeDnnQos": {"DnnMbrUplink": 100000, "DnnMbrDownlink": 50000, "TrafficClass": {"Qci": 9, "Arp": 6}}
				},
				"Imsi": ["123456789123456", "123456789123457", "123456789123458", "123456789123459", "123456789123460"]}
		 ],
		 "AppFilters": {
			"PccRuleBase": [{"FlowInfos": [{"FlowDesc": "permit out ip from 8.8.8.8/32 to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule1", "Priority": 15, "Qos": {"Var5qi": 8}},
						   {"FlowInfos": [{"FlowDesc": "permit out ip from any to assigned", "TosTrafficClass": "IPV4", "FlowDir": 2}], "RuleId": "rule2", "Priority": 25}]
		 },
		 "OperationType": 2,
		 "DeletedImsis": ["123456789123456", "123456789123457", "123456789123458", "123456789123459", "123456789123460", "123456789123461"]
		}
		]}`)

func TestUpdatePcfSubsriberPolicyDataDel(t *testing.T) {
	var nrp protos.NetworkSliceResponse
	err := json.Unmarshal(DelData, &nrp)
	if err != nil {
		panic(err)
	}
	for _, ns := range nrp.NetworkSlice {
		PCFTest.UpdatePcfSubsriberPolicyData(ns)
	}
	self := context.PCF_Self()
	if len(self.PcfSubscriberPolicyData) == 0 {
		fmt.Printf("test case TestUpdatePcfSubsriberPolicyDataDelete Passed\n")
	} else {
		t.Errorf("test case failed\n")
	}
}

func TestGetBitRateUnit(t *testing.T) {
	fmt.Printf("test case TestGetBitRateUnit \n")
	for value, expVal := range bitRateValues {
		val, unit := service.GetBitRateUnit(value)
		require.Equal(t, strconv.FormatInt(val, 10)+unit, expVal)
	}
}

func TestRegisterNF(t *testing.T) {
	// Save current function and restore at the end:
	origRegisterNFInstance := consumer.SendRegisterNFInstance
	origSearchNFInstances := consumer.SendSearchNFInstances
	origUpdateNFInstance := consumer.SendUpdateNFInstance
	defer func() {
		consumer.SendRegisterNFInstance = origRegisterNFInstance
		consumer.SendSearchNFInstances = origSearchNFInstances
		consumer.SendUpdateNFInstance = origUpdateNFInstance
	}()
	fmt.Printf("test case TestRegisterNF \n")
	var prof models.NfProfile
	consumer.SendRegisterNFInstance = func(nrfUri string, nfInstanceId string, profile models.NfProfile) (models.NfProfile, string, string, error) {
		prof = profile
		prof.HeartBeatTimer = 1
		fmt.Printf("Test RegisterNFInstance called\n")
		return prof, "", "", nil
	}
	consumer.SendSearchNFInstances = func(nrfUri string, targetNfType, requestNfType models.NfType, param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (*models.SearchResult, error) {
		fmt.Printf("Test SearchNFInstance called\n")
		return &models.SearchResult{}, nil
	}
	consumer.SendUpdateNFInstance = func(patchItem []models.PatchItem) (nfProfile models.NfProfile, problemDetails *models.ProblemDetails, err error) {
		return prof, nil, nil
	}
	go PCFTest.RegisterNF()
	service.ConfigPodTrigger <- true
	time.Sleep(5 * time.Second)
	require.Equal(t, service.KeepAliveTimer != nil, true)

	service.ConfigPodTrigger <- false
	time.Sleep(1 * time.Second)
	require.Equal(t, service.KeepAliveTimer == nil, true)
}
