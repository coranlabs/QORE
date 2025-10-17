package gmm

import (
	"fmt"
	"testing"

	"github.com/coranlabs/CORAN_LIB_UTIL/fsm"
)

func TestGmmFSM(t *testing.T) {
	if err := fsm.ExportDot(GmmFSM, "gmm"); err != nil {
		fmt.Printf("fsm export data return error: %+v", err)
	}
}
