// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func (s *Session) UpdatePDR(lSeid uint64, req *infoElement.IE) error {
	ies, err := req.UpdatePDR()
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
			// println(v)
			logger.AppLog.Trace(v)
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
				break
			}
			logger.AppLog.Trace("pdi err:", err)
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
		}
	}
	return nil
}
