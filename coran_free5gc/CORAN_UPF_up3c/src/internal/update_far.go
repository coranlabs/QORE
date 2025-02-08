// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func (s *Session) UpdateFAR(lSeid uint64, req *infoElement.IE) error {
	ies, err := req.UpdateFAR()
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
			s.fars.farID = v
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
			s.fars.applyAction = act.Flags
		case infoElement.UpdateForwardingParameters:
			xs, err := i.UpdateForwardingParameters()
			if err != nil {
				return err
			}
			v := s.newForwardingParameter(xs)
			if v != nil {
				break
			}
			logger.AppLog.Debugf("forwarding parameters %v", v)
		}
	}
	return nil
}
