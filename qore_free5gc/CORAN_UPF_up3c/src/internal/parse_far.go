// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func (s *Session) CreateFAR(lSeid uint64, req *infoElement.IE) error {
	s.fars.LocalSEID = lSeid

	ies, err := req.CreateFAR()
	if err != nil {
		return err
	}

	for _, i := range ies {
		switch i.Type {
		case infoElement.FARID:
			v, err := i.FARID()
			if err != nil {
				return err
			}
			s.fars.farID = uint32(v)
		case infoElement.ApplyAction:
			b, err := i.ApplyAction()
			if err != nil {
				return err
			}
			var act ApplyAction
			err = act.Unmarshal(b)
			if err != nil {
				return err
			}
			logger.AppLog.Debugf("apply action: %d", act.Flags)
			// Assuming ApplyAction has a method to get relevant fields
			s.fars.applyAction = act.Flags
		case infoElement.ForwardingParameters:
			xs, err := i.ForwardingParameters()
			if err != nil {
				return err
			}
			v := s.newForwardingParameter(xs)
			if v != nil {
				return err
			}
			logger.AppLog.Debugf("forwarding parameters %v", v)

		case infoElement.BARID:
			logger.AppLog.Infof("not handling barid")
		}
	}
	logger.Pfcplog.Debugf("fars extracted %v", s.fars)
	return nil
}

func (s *Session) newForwardingParameter(ies []*infoElement.IE) error {
	for _, x := range ies {
		switch x.Type {
		case infoElement.DestinationInterface:
			v, err := x.DestinationInterface()
			if err != nil {
				break
			}
			logger.AppLog.Debugf("Destination interface: %d", v)
			s.fars.forwardingparameters.DestinationInterface = v

		case infoElement.NetworkInstance:
			v, err := x.NetworkInstance()
			if err != nil {
				break
			}
			s.fars.NetworkInstance = v
		case infoElement.OuterHeaderCreation:
			v, err := x.OuterHeaderCreation()
			if err != nil {
				break
			}
			s.fars.forwardingparameters.OuterHeaderCreation.OuterHeaderCreationDescription = v.OuterHeaderCreationDescription
			if x.HasTEID() {
				s.fars.forwardingparameters.OuterHeaderCreation.TEID = v.TEID
				// GTPv1-U port
				// far.forwardingparameters.OuterHeaderCreation.Port = 2152
			} else {
				s.fars.forwardingparameters.OuterHeaderCreation.Port = v.PortNumber
			}
			if x.HasIPv4() {
				s.fars.forwardingparameters.OuterHeaderCreation.IPv4 = v.IPv4Address
			}
			logger.AppLog.Debugf("outerheadercereatin values: desc: %d , teid: %d , port: %d , ipv4: %s", v.OuterHeaderCreationDescription, v.TEID, v.PortNumber, v.IPv4Address)
		}
	}

	return nil
}
