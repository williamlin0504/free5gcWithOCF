package producer

import (
	"free5gcWithOCF/lib/http_wrapper"
	"free5gcWithOCF/lib/openapi/models"
	"free5gcWithOCF/src/ocf/context"
	"free5gcWithOCF/src/ocf/logger"
	"net/http"
	"strconv"
	"time"
)

func HandleCreateOCFEventSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	createEventSubscription := request.Body.(models.OcfCreateEventSubscription)

	createdEventSubscription, problemDetails := CreateOCFEventSubscriptionProcedure(createEventSubscription)
	if createdEventSubscription != nil {
		return http_wrapper.NewResponse(http.StatusCreated, nil, createdEventSubscription)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		return http_wrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)
	}
}

// TODO: handle event filter
func CreateOCFEventSubscriptionProcedure(createEventSubscription models.OcfCreateEventSubscription) (
	*models.OcfCreatedEventSubscription, *models.ProblemDetails) {

	ocfSelf := context.OCF_Self()

	createdEventSubscription := &models.OcfCreatedEventSubscription{}
	subscription := createEventSubscription.Subscription
	contextEventSubscription := &context.OCFContextEventSubscription{}
	contextEventSubscription.EventSubscription = *subscription
	var isImmediate bool
	var immediateFlags []bool
	var reportlist []models.OcfEventReport

	id, err := ocfSelf.EventSubscriptionIDGenerator.Allocate()
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		return nil, problemDetails
	}
	newSubscriptionID := strconv.Itoa(int(id))

	// store subscription in context
	ueEventSubscription := context.OcfUeEventSubscription{}
	ueEventSubscription.EventSubscription_OCF = &contextEventSubscription.EventSubscription
	ueEventSubscription.Timestamp = time.Now().UTC()

	if subscription.Options != nil && subscription.Options.Trigger == models.OcfEventTrigger_CONTINUOUS {
		ueEventSubscription.RemainReports = new(int32)
		*ueEventSubscription.RemainReports = subscription.Options.MaxReports
	}
	for _, eventsOCF := range *subscription.EventListOCF {
		immediateFlags = append(immediateFlags, eventsOCF.ImmediateFlag)
		if eventsOCF.ImmediateFlag {
			isImmediate = true
		}
	}

	if subscription.AnyUE {
		contextEventSubscription.IsAnyUe = true
		ueEventSubscription.AnyUe = true
		ocfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.OcfUe)
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.OcfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			return true
		})
	} else if subscription.GroupId != "" {
		contextEventSubscription.IsGroupUe = true
		ueEventSubscription.AnyUe = true
		ocfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.OcfUe)
			if ue.GroupID == subscription.GroupId {
				ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.OcfUeEventSubscription)
				*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
				contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
			}
			return true
		})

	} else {
		if ue, ok := ocfSelf.OcfUeFindBySupi(subscription.Supi); !ok {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "UE_NOT_SERVED_BY_OCF",
			}
			return nil, problemDetails
		} else {
			ue.EventSubscriptionsInfo[newSubscriptionID] = new(context.OcfUeEventSubscription)
			*ue.EventSubscriptionsInfo[newSubscriptionID] = ueEventSubscription
			contextEventSubscription.UeSupiList = append(contextEventSubscription.UeSupiList, ue.Supi)
		}
	}

	// delete subscription
	if subscription.Options != nil {
		contextEventSubscription.Expiry = subscription.Options.Expiry
	}
	ocfSelf.NewEventSubscription(newSubscriptionID, contextEventSubscription)

	// build response

	createdEventSubscription.Subscription = subscription
	createdEventSubscription.SubscriptionId = newSubscriptionID

	// for immediate use
	if subscription.AnyUE {
		ocfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.OcfUe)
			if isImmediate {
				subReports(ue, newSubscriptionID)
			}
			for i, flag := range immediateFlags {
				if flag {
					report, ok := NewOcfEventReport(ue, (*subscription.EventListOCF)[i].Type, newSubscriptionID)
					if ok {
						reportlist = append(reportlist, report)
					}
				}
			}
			// delete subscription
			if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
				delete(ue.EventSubscriptionsInfo, newSubscriptionID)
			}
			return true
		})
	} else if subscription.GroupId != "" {
		ocfSelf.UePool.Range(func(key, value interface{}) bool {
			ue := value.(*context.OcfUe)
			if isImmediate {
				subReports(ue, newSubscriptionID)
			}
			if ue.GroupID == subscription.GroupId {
				for i, flag := range immediateFlags {
					if flag {
						report, ok := NewOcfEventReport(ue, (*subscription.EventListOCF)[i].Type, newSubscriptionID)
						if ok {
							reportlist = append(reportlist, report)
						}
					}
				}
				// delete subscription
				if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
					delete(ue.EventSubscriptionsInfo, newSubscriptionID)
				}
			}
			return true
		})
	} else {
		ue, _ := ocfSelf.OcfUeFindBySupi(subscription.Supi)
		if isImmediate {
			subReports(ue, newSubscriptionID)
		}
		for i, flag := range immediateFlags {
			if flag {
				report, ok := NewOcfEventReport(ue, (*subscription.EventListOCF)[i].Type, newSubscriptionID)
				if ok {
					reportlist = append(reportlist, report)
				}
			}
		}
		// delete subscription
		if len := len(reportlist); len > 0 && (!reportlist[len-1].State.Active) {
			delete(ue.EventSubscriptionsInfo, newSubscriptionID)
		}
	}
	if len(reportlist) > 0 {
		createdEventSubscription.ReportList = reportlist
		// delete subscription
		if !reportlist[0].State.Active {
			ocfSelf.DeleteEventSubscription(newSubscriptionID)
		}
	}

	return createdEventSubscription, nil
}

func HandleDeleteOCFEventSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	logger.EeLog.Infoln("Handle Delete OCF Event Subscription")

	subscriptionID := request.Params["subscriptionId"]

	problemDetails := DeleteOCFEventSubscriptionProcedure(subscriptionID)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, nil)
	}
}

func DeleteOCFEventSubscriptionProcedure(subscriptionID string) *models.ProblemDetails {
	ocfSelf := context.OCF_Self()

	subscription, ok := ocfSelf.FindEventSubscription(subscriptionID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "SUBSCRIPTION_NOT_FOUND",
		}
		return problemDetails
	}

	for _, supi := range subscription.UeSupiList {
		if ue, ok := ocfSelf.OcfUeFindBySupi(supi); ok {
			delete(ue.EventSubscriptionsInfo, subscriptionID)
		}
	}
	ocfSelf.DeleteEventSubscription(subscriptionID)
	return nil
}

func HandleModifyOCFEventSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	logger.EeLog.Infoln("Handle Modify OCF Event Subscription")

	subscriptionID := request.Params["subscriptionId"]
	modifySubscriptionRequest := request.Body.(models.ModifySubscriptionRequest)

	updatedEventSubscription, problemDetails := ModifyOCFEventSubscriptionProcedure(subscriptionID,
		modifySubscriptionRequest)
	if updatedEventSubscription != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, updatedEventSubscription)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		return http_wrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)
	}
}

func ModifyOCFEventSubscriptionProcedure(
	subscriptionID string,
	modifySubscriptionRequest models.ModifySubscriptionRequest) (
	*models.OcfUpdatedEventSubscription, *models.ProblemDetails) {

	ocfSelf := context.OCF_Self()

	contextSubscription, ok := ocfSelf.FindEventSubscription(subscriptionID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "SUBSCRIPTION_NOT_FOUND",
		}
		return nil, problemDetails
	}

	if modifySubscriptionRequest.OptionItem != nil {
		contextSubscription.Expiry = modifySubscriptionRequest.OptionItem.Value
	} else if modifySubscriptionRequest.SubscriptionItemInner != nil {
		subscription := &contextSubscription.EventSubscription
		if !contextSubscription.IsAnyUe && !contextSubscription.IsGroupUe {
			if _, ok := ocfSelf.OcfUeFindBySupi(subscription.Supi); !ok {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusForbidden,
					Cause:  "UE_NOT_SERVED_BY_OCF",
				}
				return nil, problemDetails
			}
		}
		op := modifySubscriptionRequest.SubscriptionItemInner.Op
		index, err := strconv.Atoi(modifySubscriptionRequest.SubscriptionItemInner.Path[11:])
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			return nil, problemDetails
		}
		lists := (*subscription.EventListOCF)
		len := len(*subscription.EventListOCF)
		switch op {
		case "replace":
			eventOCF := *modifySubscriptionRequest.SubscriptionItemInnerOCF.Value
			if index < len {
				(*subscription.EventListOCF)[index] = eventOCF
			}
		case "remove":
			if index < len {
				*subscription.EventListOCF = append(lists[:index], lists[index+1:]...)
			}
		case "add":
			eventOCF := *modifySubscriptionRequest.SubscriptionItemInnerOCF.Value
			*subscription.EventListOCF = append(lists, eventOCF)
		}
	}

	updatedEventSubscription := &models.OcfUpdatedEventSubscription{
		Subscription: &contextSubscription.EventSubscription,
	}
	return updatedEventSubscription, nil
}

func subReports(ue *context.OcfUe, subscriptionId string) {
	remainReport := ue.EventSubscriptionsInfo[subscriptionId].RemainReports
	if remainReport == nil {
		return
	}
	*remainReport--
}

// DO NOT handle OcfEventType_PRESENCE_IN_AOI_REPORT and OcfEventType_UES_IN_AREA_REPORT(about area)
func NewOcfEventReport(ue *context.OcfUe, Type models.OcfEventType, subscriptionId string) (
	report models.OcfEventReport, ok bool) {
	ueSubscription, ok := ue.EventSubscriptionsInfo[subscriptionId]
	if !ok {
		return report, ok
	}

	report.AnyUe = ueSubscription.AnyUe
	report.Supi = ue.Supi
	report.Type = Type
	report.TimeStamp = &ueSubscription.Timestamp
	report.State = new(models.OcfEventState)
	mode := ueSubscription.EventSubscription_OCF.Options
	if mode == nil {
		report.State.Active = true
	} else if mode.Trigger == models.OcfEventTrigger_ONE_TIME {
		report.State.Active = false
	} else if *ueSubscription.RemainReports <= 0 {
		report.State.Active = false
	} else {
		report.State.Active = getDuration(mode.Expiry, &report.State.RemainDuration)
		if report.State.Active {
			report.State.RemainReports = *ueSubscription.RemainReports
		}
	}

	switch Type {
	case models.OcfEventType_LOCATION_REPORT:
		report.Location = &ue.Location
	// case models.OcfEventType_PRESENCE_IN_AOI_REPORT:
	// report.AreaList = (*subscription.EventListOCF)[eventIndex].AreaList
	case models.OcfEventType_TIMEZONE_REPORT:
		report.Timezone = ue.TimeZone
	case models.OcfEventType_ACCESS_TYPE_REPORT:
		for accessType, state := range ue.State {
			if state.Is(context.Registered) {
				report.AccessTypeList = append(report.AccessTypeList, accessType)
			}
		}
	case models.OcfEventType_REGISTRATION_STATE_REPORT:
		var rmInfos []models.RmInfo
		for accessType, state := range ue.State {
			rmInfo := models.RmInfo{
				RmState:    models.RmState_DEREGISTERED,
				AccessType: accessType,
			}
			if state.Is(context.Registered) {
				rmInfo.RmState = models.RmState_REGISTERED
			}
			rmInfos = append(rmInfos, rmInfo)
		}
		report.RmInfoList = rmInfos
	case models.OcfEventType_CONNECTIVITY_STATE_REPORT:
		report.CmInfoList = ue.GetCmInfo()
	case models.OcfEventType_REACHABILITY_REPORT:
		report.Reachability = ue.Reachability
	case models.OcfEventType_SUBSCRIBED_DATA_REPORT:
		report.SubscribedData = &ue.SubscribedData
	case models.OcfEventType_COMMUNICATION_FAILURE_REPORT:
		// TODO : report.CommFailure
	case models.OcfEventType_SUBSCRIPTION_ID_CHANGE:
		report.SubscriptionId = subscriptionId
	case models.OcfEventType_SUBSCRIPTION_ID_ADDITION:
		report.SubscriptionId = subscriptionId
	}
	return report, ok

}

func getDuration(expiry *time.Time, remainDuration *int32) bool {

	if expiry != nil {
		if time.Now().After(*expiry) {
			return false
		} else {
			duration := time.Until(*expiry)
			*remainDuration = int32(duration.Seconds())
		}
	}
	return true

}
