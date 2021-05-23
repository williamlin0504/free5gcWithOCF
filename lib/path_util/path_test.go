package path_util

import (
	" free5gcWithOCF/lib/path_util/logger"
	"testing"
)

func TestFree5gcPath(t *testing.T) {
	logger.PathLog.Infoln(Go free5gcWithOCFPath(" free5gcWithOCF/abcdef/abcdef.pem"))
}
