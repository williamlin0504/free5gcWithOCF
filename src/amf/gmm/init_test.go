package gmm_test

import (
	"free5gcWithOCF/lib/fsm"
	"free5gcWithOCF/src/amf/gmm"
	"testing"
)

func TestGmmFSM(t *testing.T) {
	fsm.ExportDot(gmm.GmmFSM, "gmm")
}
