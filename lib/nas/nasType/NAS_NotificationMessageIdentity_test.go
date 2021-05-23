package nasType_test

import (
	" free5gc/lib/nas"
	" free5gcs/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNotificationMessageIdentity(t *testing.T) {
	a := nasType.NewNotificationMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypeNotificationMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypeNotificationMessageIdentityMessageTypeTable = []nasTypeNotificationMessageIdentityMessageType{
	{nas.MsgTypeNotification, nas.MsgTypeNotification},
}

func TestNasTypeGetSetNotificationMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewNotificationMessageIdentity()
	for _, table := range nasTypeNotificationMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
