package nas

import (
	"fmt"

	"github.com/free5gc/ocf/context"
	"github.com/free5gc/ocf/logger"
	"github.com/free5gc/ocf/nas/nas_security"
)

func HandleNAS(ue *context.RanUe, procedureCode int64, nasPdu []byte) {
	ocfSelf := context.OCF_Self()

	if ue == nil {
		logger.NasLog.Error("RanUe is nil")
		return
	}

	if nasPdu == nil {
		ue.Log.Error("nasPdu is nil")
		return
	}

	if ue.OcfUe == nil {
		ue.OcfUe = ocfSelf.NewOcfUe("")
		ue.OcfUe.AttachRanUe(ue)

		// set log information
		ue.OcfUe.NASLog = logger.NasLog.WithField(logger.FieldOcfUeNgapID, fmt.Sprintf("OCF_UE_NGAP_ID:%d", ue.OcfUeNgapId))
		ue.OcfUe.GmmLog = logger.GmmLog.WithField(logger.FieldOcfUeNgapID, fmt.Sprintf("OCF_UE_NGAP_ID:%d", ue.OcfUeNgapId))
	}

	msg, err := nas_security.Decode(ue.OcfUe, ue.Ran.AnType, nasPdu)
	if err != nil {
		ue.OcfUe.NASLog.Errorln(err)
		return
	}

	if err := Dispatch(ue.OcfUe, ue.Ran.AnType, procedureCode, msg); err != nil {
		ue.OcfUe.NASLog.Errorf("Handle NAS Error: %v", err)
	}
}
