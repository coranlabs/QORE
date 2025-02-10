<!--
SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
Copyright 2019 free5GC.org

SPDX-License-Identifier: Apache-2.0

-->

# pcf

The Policy Control Function (PCF) supports unified policy framework to govern network behaviour, 
provides policy rules to Control Plane function(s) to enforce them and Accesses subscription information
relevant for policy decisions in a Unified Data Repository (UDR)

The reference 3GPP specification for PCF are as follows
- PCC framework Specification 23.503,
- Session Management Policy Control - Specification 29.512
- Policy and Charging Control signalling flows and QoS parameter mapping- Specification 29513


## PCF Block Diagram
![PCF Block Diagram](/docs/images/README-PCF.png)

## Supported Features
- PCF provides Access and Mobility Management related policies to the AMF Subscription Data retrieval and AM Policy management
- PCF provides Session Management Policy Control Service to the SMF Subscription Data retrieval and SM Policy management
- Policy Control Function (PCF) shall support interactions with the access and mobility policy enforcement in the AMF, through service-based interfaces

## Upcoming Changes in PCF
- Process configuration received from Configuration Service and prepare PCC Rules, Session Rules and Qos Flows 
- Send PCC Rules, Session Rules to SMF when a SMF Creates Policy subscriber
- Send notification towards SMF PDU Session when PCF detects any changes in Subscriber’s Rule/Qos information.
- Dedicated QoS flows addition  & removal through APIs



Compliance of the 5G Network functions can be found at [5G Compliance ](https://docs.sd-core.opennetworking.org/master/overview/3gpp-compliance-5g.html)

## Reach out to us thorugh 

1. #sdcore-dev channel in [ONF Community Slack](https://onf-community.slack.com/)
2. Raise Github issues
