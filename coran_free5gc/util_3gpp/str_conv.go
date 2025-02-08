// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package util_3gpp

import (
	"fmt"

	"github.com/omec-project/openapi/models"
)

func SNssaiToString(snssai *models.Snssai) (str string) {
	if snssai.Sd == "" {
		return fmt.Sprintf("%d-%s", snssai.Sst, snssai.Sd)
	}
	return fmt.Sprintf("%d", snssai.Sst)
}
