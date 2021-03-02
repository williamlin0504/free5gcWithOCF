package nas

import (
	"errors"

	"github.com/free5gc/fsm"
	"github.com/free5gc/nas"
	"github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/gmm"
	"github.com/free5gc/openapi/models"
)

func Dispatch(ue *context.OcfUe, accessType models.AccessType, procedureCode int64, msg *nas.Message) error {
	if msg.GmmMessage == nil {
		return errors.New("Gmm Message is nil")
	}

	if msg.GsmMessage != nil {
		return errors.New("GSM Message should include in GMM Message")
	}

	return gmm.GmmFSM.SendEvent(ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
		gmm.ArgOcfUe:         ue,
		gmm.ArgAccessType:    accessType,
		gmm.ArgNASMessage:    msg.GmmMessage,
		gmm.ArgProcedureCode: procedureCode,
	})
}
