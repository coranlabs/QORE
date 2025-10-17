// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	"github.com/coranlabs/HEXA_UPF/src/logger"
	infoElement "github.com/wmnsk/go-pfcp/ie"
)

func (s *Session) CreatePDR(lSeid uint64, req *infoElement.IE) error {

	s.pdrs.LocalSEID = lSeid
	ies, err := req.CreatePDR()
	if err != nil {
		return err
	}

	for _, i := range ies {
		switch i.Type {
		case infoElement.PDRID:
			v, err := i.PDRID()
			if err != nil {
				break
			}
			s.pdrs.PDRID = uint32(v)

		case infoElement.Precedence:
			v, err := i.Precedence()
			if err != nil {
				break
			}

			s.pdrs.Precedence = v

		case infoElement.PDI:
			err := s.newPdi(i)
			if err != nil {
				logger.AppLog.Fatalln("no pdi called error ", err)
			}
		case infoElement.OuterHeaderRemoval:
			v, err := i.OuterHeaderRemovalDescription()
			if err != nil {
				break
			}

			s.pdrs.outerHeaderRemoval = v
		case infoElement.FARID:
			v, err := i.FARID()
			if err != nil {
				break
			}

			s.pdrs.FARID = v
		case infoElement.QERID:
			v, err := i.QERID()
			if err != nil {
				break
			}
			s.pdrs.QERID = v
		}
	}

	logger.AppLog.Debug("create pdr function called")
	logger.AppLog.Debugf("pdrs extracted: %v", s.pdrs)
	return nil
}

func (s *Session) newPdi(i *infoElement.IE) error {
	ies, err := i.PDI()
	if err != nil {
		return err
	}
	for _, x := range ies {
		switch x.Type {
		case infoElement.SourceInterface:
			v, err := x.SourceInterface()
			if err != nil {
				break
			}
			s.pdrs.PDI.SourceInterface = v
		case infoElement.FTEID:
			v, err := x.FTEID()
			if err != nil {
				break
			}
			s.pdrs.FTEID.TEID = v.TEID
			s.pdrs.FTEID.IPv4Address = v.IPv4Address
		case infoElement.NetworkInstance:
			v, err := x.NetworkInstance()
			if err != nil {
				break
			}
			s.pdrs.PDI.NetworkInstance = v
		case infoElement.UEIPAddress:
			v, err := x.UEIPAddress()
			if err != nil {
				break
			}
			s.pdrs.PDI.UeIpAddress = v.IPv4Address
		case infoElement.SDFFilter:
			logger.AppLog.Infof("not handling SDFFilter for now")
		case infoElement.ApplicationID:
			logger.AppLog.Infof("not handling Application id for now")
		}
	}
	return nil
}
