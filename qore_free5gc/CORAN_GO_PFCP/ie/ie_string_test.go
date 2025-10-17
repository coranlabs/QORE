// Copyright 2019-2024 go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"testing"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/google/go-cmp/cmp"
)

func TestStringIEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		decoded     string
		decoderFunc func(*ie.IE) (string, error)
	}{
		{
			description: "ActivatePredefinedRules",
			structured:  ie.NewActivatePredefinedRules("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ActivatePredefinedRules() },
		}, {
			description: "APNDNN",
			structured:  ie.NewAPNDNN("some.apn.example"),
			decoded:     "some.apn.example",
			decoderFunc: func(i *ie.IE) (string, error) { return i.APNDNN() },
		}, {
			description: "ApplicationID",
			structured:  ie.NewApplicationID("https://github.com/wmnsk/go-pfcp/"),
			decoded:     "https://github.com/wmnsk/go-pfcp/",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ApplicationID() },
		}, {
			description: "ApplicationInstanceID",
			structured:  ie.NewApplicationInstanceID("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ApplicationInstanceID() },
		}, {
			description: "BridgeManagementInformationContainer",
			structured:  ie.NewBridgeManagementInformationContainer("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.BridgeManagementInformationContainer() },
		}, {
			description: "BridgeManagementInformationContainer/TSCManagementInformationWithinSessionModificationRequest",
			structured: ie.NewTSCManagementInformationWithinSessionModificationRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.BridgeManagementInformationContainer() },
		}, {
			description: "BridgeManagementInformationContainer/TSCManagementInformationWithinSessionModificationResponse",
			structured: ie.NewTSCManagementInformationWithinSessionModificationResponse(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.BridgeManagementInformationContainer() },
		}, {
			description: "BridgeManagementInformationContainer/TSCManagementInformationWithinSessionReportRequest",
			structured: ie.NewTSCManagementInformationWithinSessionReportRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.BridgeManagementInformationContainer() },
		}, {
			description: "DataNetworkAccessIdentifier",
			structured:  ie.NewDataNetworkAccessIdentifier("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.DataNetworkAccessIdentifier() },
		}, {
			description: "DeactivatePredefinedRules",
			structured:  ie.NewDeactivatePredefinedRules("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.DeactivatePredefinedRules() },
		}, {
			description: "FlowInformation/FlowDescription",
			structured:  ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FlowDescription() },
		}, {
			description: "ForwardingPolicyIdentifier",
			structured:  ie.NewForwardingPolicy("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ForwardingPolicyIdentifier() },
		}, {
			description: "FramedIPv6Route",
			structured:  ie.NewFramedIPv6Route("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FramedIPv6Route() },
		}, {
			description: "FramedRoute",
			structured:  ie.NewFramedRoute("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FramedRoute() },
		}, {
			description: "NetworkInstance",
			structured:  ie.NewNetworkInstance("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstance() },
		}, {
			description: "NetworkInstance/Heuristic",
			structured:  ie.NewNetworkInstance("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceHeuristic() },
		}, {
			description: "NetworkInstance/FQDN",
			structured:  ie.NewNetworkInstanceFQDN("go-pfcp.com"),
			decoded:     "go-pfcp.com",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceFQDN() },
		}, {
			description: "NetworkInstance/FQDN/Heuristic",
			structured:  ie.NewNetworkInstanceFQDN("go-pfcp.com"),
			decoded:     "go-pfcp.com",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceHeuristic() },
		}, {
			description: "NodeID/IPv4",
			structured:  ie.NewNodeID("127.0.0.1", "", ""),
			decoded:     "127.0.0.1",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "NodeID/IPv6",
			structured:  ie.NewNodeID("", "2001::1", ""),
			decoded:     "2001::1",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "NodeID/FQDN",
			structured:  ie.NewNodeID("", "", "go-pfcp.epc.3gppnetwork.org"),
			decoded:     "go-pfcp.epc.3gppnetwork.org",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "NodeIDHeuristic/IPv4",
			structured:  ie.NewNodeIDHeuristic("127.0.0.1"),
			decoded:     "127.0.0.1",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "NodeIDHeuristic/IPv6",
			structured:  ie.NewNodeIDHeuristic("2001::1"),
			decoded:     "2001::1",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "NodeIDHeuristic/FQDN",
			structured:  ie.NewNodeIDHeuristic("go-pfcp.epc.3gppnetwork.org"),
			decoded:     "go-pfcp.epc.3gppnetwork.org",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NodeID() },
		}, {
			description: "PortManagementInformationContainer",
			structured:  ie.NewPortManagementInformationContainer("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.PortManagementInformationContainer() },
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionModificationRequest",
			structured: ie.NewTSCManagementInformationWithinSessionModificationRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.PortManagementInformationContainer() },
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionModificationResponse",
			structured: ie.NewTSCManagementInformationWithinSessionModificationResponse(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.PortManagementInformationContainer() },
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionReportRequest",
			structured: ie.NewTSCManagementInformationWithinSessionReportRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewBridgeManagementInformationContainer("go-pfcp"),
			),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.PortManagementInformationContainer() },
		}, {
			description: "SMFSetID",
			structured:  ie.NewSMFSetID("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.SMFSetID() },
		}, {
			description: "UEIPAddressPoolIdentity",
			structured:  ie.NewUEIPAddressPoolIdentity("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.UEIPAddressPoolIdentityString() },
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			got, err := c.decoderFunc(c.structured)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, c.decoded); diff != "" {
				t.Error(diff)
			}
		})
	}
}
