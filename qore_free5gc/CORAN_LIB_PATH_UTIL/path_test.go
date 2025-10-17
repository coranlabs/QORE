package path_util

import (
	"testing"

	"github.com/coranlabs/CORAN_LIB_PATH_UTIL/logger"
)

func TestCoranlabsPath(t *testing.T) {
	logger.PathLog.Infoln(CoranlabsPath("coranlabs/abcdef/abcdef.pem"))
}
