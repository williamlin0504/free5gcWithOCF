package path_util

import (
	" free5gc/lib/path_util/logger"
	"testing"
)

func TestFree5gcPath(t *testing.T) {
	logger.PathLog.Infoln(Go free5gcPath(" free5gc/abcdef.pem"))
}
