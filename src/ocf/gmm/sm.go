package gmm

import (
	"free5gc/lib/fsm"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
	"free5gc/src/ocf/context"
	gmm_message "free5gc/src/ocf/gmm/message"
	"free5gc/src/ocf/logger"
)

func DeRegistered(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State[DeRegistered]")

	switch event {
	case fsm.EntryEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		ocfUe.ClearRegistrationRequestData(accessType)
	case GmmMessageEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeRegistrationRequest:
			if err := HandleRegistrationRequest(ocfUe, accessType, procedureCode, gmmMessage.RegistrationRequest); err != nil {
				logger.GmmLog.Errorln(err)
			} else {
				if err := GmmFSM.SendEvent(state, StartAuthEvent, fsm.ArgsType{
					ArgOcfUe:         ocfUe,
					ArgAccessType:    accessType,
					ArgProcedureCode: procedureCode,
				}); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case StartAuthEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func Registered(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State[Registered]")

	switch event {
	case fsm.EntryEvent:
		// clear stored registration request data for this registration
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		ocfUe.ClearRegistrationRequestData(accessType)
	case GmmMessageEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		switch gmmMessage.GetMessageType() {
		// Mobility Registration update / Periodic Registration update
		case nas.MsgTypeRegistrationRequest:
			if err := HandleRegistrationRequest(ocfUe, accessType, procedureCode, gmmMessage.RegistrationRequest); err != nil {
				logger.GmmLog.Errorln(err)
			} else {
				if err := GmmFSM.SendEvent(state, StartAuthEvent, fsm.ArgsType{
					ArgOcfUe:         ocfUe,
					ArgAccessType:    accessType,
					ArgProcedureCode: procedureCode,
				}); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		case nas.MsgTypeULNASTransport:
			if err := HandleULNASTransport(ocfUe, accessType, gmmMessage.ULNASTransport); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeConfigurationUpdateComplete:
			if err := HandleConfigurationUpdateComplete(ocfUe, gmmMessage.ConfigurationUpdateComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeServiceRequest:
			if err := HandleServiceRequest(ocfUe, accessType, gmmMessage.ServiceRequest); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeNotificationResponse:
			if err := HandleNotificationResponse(ocfUe, gmmMessage.NotificationResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
			if err := GmmFSM.SendEvent(state, InitDeregistrationEvent, fsm.ArgsType{
				ArgOcfUe:      ocfUe,
				ArgAccessType: accessType,
				ArgNASMessage: gmmMessage,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(ocfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case StartAuthEvent:
		logger.GmmLog.Debugln(event)
	case InitDeregistrationEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func Authentication(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State [Authentication]")

	switch event {
	case fsm.EntryEvent:
		fallthrough
	case AuthRestartEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		accessType := args[ArgAccessType].(models.AccessType)

		pass, err := AuthenticationProcedure(ocfUe, accessType)
		if err != nil {
			logger.GmmLog.Errorln(err)
		}
		if pass {
			if err := GmmFSM.SendEvent(state, AuthSuccessEvent, fsm.ArgsType{
				ArgOcfUe:      ocfUe,
				ArgAccessType: accessType,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		}
	case GmmMessageEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)

		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeIdentityResponse:
			if err := HandleIdentityResponse(ocfUe, gmmMessage.IdentityResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
			err := GmmFSM.SendEvent(state, AuthRestartEvent, fsm.ArgsType{ArgOcfUe: ocfUe, ArgAccessType: accessType})
			if err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeAuthenticationResponse:
			if err := HandleAuthenticationResponse(ocfUe, accessType, gmmMessage.AuthenticationResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeAuthenticationFailure:
			if err := HandleAuthenticationFailure(ocfUe, accessType, gmmMessage.AuthenticationFailure); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(ocfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case AuthSuccessEvent:
		logger.GmmLog.Debugln(event)
	case AuthFailEvent:
		logger.GmmLog.Debugln(event)
		logger.GmmLog.Warnln("Reject authentication")
	case fsm.ExitEvent:
		// clear authentication related data at exit
		logger.GmmLog.Debugln(event)
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		ocfUe.AuthenticationCtx = nil
		ocfUe.AuthFailureCauseSynchFailureTimes = 0
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func SecurityMode(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State[SecurityMode]")

	switch event {
	case fsm.EntryEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		if ocfUe.SecurityContextIsValid() {
			logger.GmmLog.Debugln("UE has a valid security context - skip security mode control procedure")
			if err := GmmFSM.SendEvent(state, SecurityModeSuccessEvent, fsm.ArgsType{
				ArgOcfUe:      ocfUe,
				ArgAccessType: accessType,
				ArgNASMessage: ocfUe.RegistrationRequest,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		} else {
			eapSuccess := args[ArgEAPSuccess].(bool)
			eapMessage := args[ArgEAPMessage].(string)
			// Select enc/int algorithm based on ue security capability & ocf's policy,
			ocfSelf := context.OCF_Self()
			ocfUe.SelectSecurityAlg(ocfSelf.SecurityAlgorithm.IntegrityOrder, ocfSelf.SecurityAlgorithm.CipheringOrder)
			// Generate KnasEnc, KnasInt
			ocfUe.DerivateAlgKey()
			gmm_message.SendSecurityModeCommand(ocfUe.RanUe[accessType], eapSuccess, eapMessage)
		}
	case GmmMessageEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeSecurityModeComplete:
			if err := HandleSecurityModeComplete(ocfUe, accessType, procedureCode, gmmMessage.SecurityModeComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeSecurityModeReject:
			if err := HandleSecurityModeReject(ocfUe, accessType, gmmMessage.SecurityModeReject); err != nil {
				logger.GmmLog.Errorln(err)
			}
			err := GmmFSM.SendEvent(state, SecurityModeFailEvent, fsm.ArgsType{
				ArgOcfUe:      ocfUe,
				ArgAccessType: accessType,
			})
			if err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(ocfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case SecurityModeSuccessEvent:
		logger.GmmLog.Debugln(event)
	case SecurityModeFailEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
		return
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func ContextSetup(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State[ContextSetup]")

	switch event {
	case fsm.EntryEvent:
		logger.GmmLog.Debugln("EntryEvent")
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		gmmMessage := args[ArgNASMessage]
		accessType := args[ArgAccessType].(models.AccessType)

		switch message := gmmMessage.(type) {
		case *nasMessage.RegistrationRequest:
			ocfUe.RegistrationRequest = message
			switch ocfUe.RegistrationType5GS {
			case nasMessage.RegistrationType5GSInitialRegistration:
				if err := HandleInitialRegistration(ocfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
				}
			case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
				fallthrough
			case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
				if err := HandleMobilityAndPeriodicRegistrationUpdating(ocfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		case *nasMessage.ServiceRequest:
			if err := HandleServiceRequest(ocfUe, accessType, message); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve wrong gmm message")
		}
	case GmmMessageEvent:
		logger.GmmLog.Debugln("GmmMessageEvent")
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeIdentityResponse:
			if err := HandleIdentityResponse(ocfUe, gmmMessage.IdentityResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
			switch ocfUe.RegistrationType5GS {
			case nasMessage.RegistrationType5GSInitialRegistration:
				if err := HandleInitialRegistration(ocfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
					err = GmmFSM.SendEvent(state, ContextSetupFailEvent, fsm.ArgsType{
						ArgOcfUe:      ocfUe,
						ArgAccessType: accessType,
					})
					if err != nil {
						logger.GmmLog.Errorln(err)
					}
				}
			case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
				fallthrough
			case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
				if err := HandleMobilityAndPeriodicRegistrationUpdating(ocfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
					err = GmmFSM.SendEvent(state, ContextSetupFailEvent, fsm.ArgsType{
						ArgOcfUe:      ocfUe,
						ArgAccessType: accessType,
					})
					if err != nil {
						logger.GmmLog.Errorln(err)
					}
				}
			}
		case nas.MsgTypeRegistrationComplete:
			if err := HandleRegistrationComplete(ocfUe, accessType, gmmMessage.RegistrationComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(ocfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorln("UE state mismatch")
		}
	case ContextSetupSuccessEvent:
		logger.GmmLog.Debugln(event)
	case ContextSetupFailEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func DeregisteredInitiated(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	logger.GmmLog.Debugln("UE in GMM State[DeregisteredInitiated]")

	switch event {
	case fsm.EntryEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		if err := HandleDeregistrationRequest(ocfUe, accessType,
			gmmMessage.DeregistrationRequestUEOriginatingDeregistration); err != nil {
			logger.GmmLog.Errorln(err)
		}
	case GmmMessageEvent:
		ocfUe := args[ArgOcfUe].(*context.OcfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
			if err := HandleDeregistrationAccept(ocfUe, accessType,
				gmmMessage.DeregistrationAcceptUETerminatedDeregistration); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case DeregistrationAcceptEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}
