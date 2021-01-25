package nas

import (
	"errors"
	"free5gcWithOCF/lib/fsm"
	"free5gcWithOCF/lib/nas"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/gmm"
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
