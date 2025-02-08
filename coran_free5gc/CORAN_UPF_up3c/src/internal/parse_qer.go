// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 CORAN LABS

package internal

import (
	infoElement "github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/src/logger"
)

func (s *Session) CreateQER(lSeid uint64, req *infoElement.IE) error {
	logger.AppLog.Infof("Create QER")
	return nil
}
