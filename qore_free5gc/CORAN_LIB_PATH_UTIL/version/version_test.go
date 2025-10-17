package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coranlabs/CORAN_LIB_PATH_UTIL/version"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "2020-03-31-01", version.GetVersion())
}
