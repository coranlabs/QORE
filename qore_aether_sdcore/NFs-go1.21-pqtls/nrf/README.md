<!--
SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
Copyright 2019 free5GC.org

SPDX-License-Identifier: Apache-2.0

-->

# NRF

Network Repository Function provides service discovery functionality in the 5G core network. 
Each network function registers with NRF with a set of properties called as NF Profile.
Implements 3gpp specification 29.510. NRF Keeps the profile data for each network function in the 
MongoDB. Supports Discovery & registration procedure. When the discovery API is called, NRF 
fetches a matching profile from the database and returns it to the caller. 


## NRF block diagram
![UDM Block Diagram](/docs/images/README-NRF.png)

## Supported Features
- Registration of Network Functions
- Searching of matching Network functions
- Handling multiple instances of registration from Network Functions
- Supporting keepalive functionality to check the health of network functions


## Upcoming changes in NRF
- Supporting callbacks to send notification when a network function is added/removed/modified.
- Subscription management callbacks to network functions.
- NRF cache library which can be used by modules to avoid frequent queries to NRF

Compliance of the 5G Network functions can be found at [5G Compliance ](https://docs.sd-core.opennetworking.org/master/overview/3gpp-compliance-5g.html)

## Reach out to us thorugh 

1. #sdcore-dev channel in [ONF Community Slack](https://onf-community.slack.com/)
2. Raise Github issues
