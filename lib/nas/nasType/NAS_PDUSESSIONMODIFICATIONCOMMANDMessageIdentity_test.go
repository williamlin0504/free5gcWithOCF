package nasType_test

import (
	" free5gcWithOCF/lib/nas"
	" free5gcWithOCF/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONMODIFICATIONCOMMANDMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONMODIFICATIONCOMMANDMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONMODIFICATIONCOMMANDMessageIdentityMessageTypeTable = []nasTypePDUSESSIONMODIFICATIONCOMMANDMessageIdentityMessageType{
	{nas.MsgTypePDUSessionModificationCommand, nas.MsgTypePDUSessionModificationCommand},
}

func TestNasTypeGetSetPDUSESSIONMODIFICATIONCOMMANDMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDMessageIdentity()
	for _, table := range nasTypePDUSESSIONMODIFICATIONCOMMANDMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
