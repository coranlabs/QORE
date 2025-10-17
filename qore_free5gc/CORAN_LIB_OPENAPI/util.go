package CORAN_LIB_OPENAPI

import (
	"strings"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

func SnssaiEqualFold(s, t models.Snssai) bool {
	if s.Sst == t.Sst && strings.EqualFold(s.Sd, t.Sd) {
		return true
	}
	return false
}
