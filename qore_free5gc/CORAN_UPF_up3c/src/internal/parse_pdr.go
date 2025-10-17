// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func (s *Session) CreatePDR(lSeid uint64, req *infoElement.IE) error {

	s.pdrs.LocalSEID = lSeid
	//2")
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
			//3")
			s.pdrs.PDRID = uint32(v)
		case infoElement.UserPlaneIPResourceInformation:
			v, err := i.PDRID()
			if err != nil {
				break
			}
			//3")
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
			//6")

			s.pdrs.outerHeaderRemoval = v
		case infoElement.FARID:
			v, err := i.FARID()
			if err != nil {
				break
			}
			//7")

			s.pdrs.FARID = v
		case infoElement.QERID:
			v, err := i.QERID()
			if err != nil {
				break
			}
			s.pdrs.QERID = v
		}
	}

	// // Store extracted values in variables
	// // pdrid, precedence, pdi, outerHeaderRemoval, farid, qerid, urrid
	logger.AppLog.Debug("create pdr function called")
	logger.AppLog.Debugf("pdrs extracted: %v", s.pdrs)
	return nil
}

func (s *Session) newPdi(i *infoElement.IE) error {
	ies, err := i.PDI()
	if err != nil {
		return err
	}
	//PDI")
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
			// logger.MainLog.Tracef("not handling network instance for now %s", v)
		case infoElement.UEIPAddress:
			v, err := x.UEIPAddress()
			if err != nil {
				break
			}
			// logger.MainLog.Tracef("not handling ueipaddress for now %d", v)
			s.pdrs.PDI.UeIpAddress = v.IPv4Address
			// println("ye print ho rha hai")
		case infoElement.SDFFilter:
			// sdfIEs = append(sdfIEs, x)
			logger.AppLog.Infof("not handling SDFFilter for now")
		case infoElement.ApplicationID:
			// Handle ApplicationID if necessary
			logger.AppLog.Infof("not handling Application id for now")
		}
	}
	return nil
}
