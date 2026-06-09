package service

import (
	"errors"
	"fmt"
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"
)

var (
	ErrVehicleOccupied    = errors.New("车辆已被占用")
	ErrVehicleNotFound    = errors.New("车辆不存在")
	ErrDispatchNotFound   = errors.New("调度记录不存在")
	ErrInvalidStatusTransition = errors.New("无效的状态转换")
	ErrQueueJumpReasonRequired = errors.New("VIP 插队必须填写原因")
)

func CreateDispatch(req *CreateDispatchRequest) (*models.Dispatch, error) {
	if req.IsVIP && req.QueueJumpReason == "" {
		return nil, ErrQueueJumpReasonRequired
	}

	now := time.Now()
	dispatchID := fmt.Sprintf("D%d", now.UnixNano())

	caller, err := getOrCreateCaller(req.CallerName, req.CallerPhone, req.IsVIP)
	if err != nil {
		return nil, err
	}

	priority := 0
	if req.IsVIP {
		priority = 100
	}

	dispatch := &models.Dispatch{
		ID:                 dispatchID,
		CallerID:           caller.ID,
		CallerName:         req.CallerName,
		CallerPhone:        req.CallerPhone,
		IsVIP:              req.IsVIP,
		QueueJumpReason:    req.QueueJumpReason,
		IllnessType:        req.IllnessType,
		IllnessDescription: req.IllnessDescription,
		Address:            req.Address,
		Longitude:          req.Longitude,
		Latitude:           req.Latitude,
		Status:             models.DispatchStatusQueued,
		CallTime:           now,
		ChangeLog: []models.ChangeLogEntry{
			{
				Timestamp: now,
				Action:    "created",
				Operator:  req.Operator,
				Reason:    "接警入队",
			},
		},
		Priority: priority,
	}

	store.GlobalStore.AddDispatch(dispatch)
	store.GlobalStore.AddToQueue(dispatchID, req.IsVIP)

	return dispatch, nil
}

type CreateDispatchRequest struct {
	CallerName         string
	CallerPhone        string
	IsVIP              bool
	QueueJumpReason    string
	IllnessType        models.IllnessType
	IllnessDescription string
	Address            string
	Longitude          float64
	Latitude           float64
	Operator           string
}

func getOrCreateCaller(name, phone string, isVIP bool) (*models.Caller, error) {
	caller, ok := store.GlobalStore.GetCallerByPhone(phone)
	if ok {
		if isVIP && !caller.IsVIP {
			caller.IsVIP = true
			store.GlobalStore.UpdateCaller(caller)
		}
		return caller, nil
	}

	now := time.Now()
	callerID := fmt.Sprintf("C%d", now.UnixNano())
	caller = &models.Caller{
		ID:          callerID,
		Name:        name,
		Phone:       phone,
		IsVIP:       isVIP,
		DispatchIDs: []string{},
		CreatedAt:   now,
	}
	store.GlobalStore.AddCaller(caller)
	return caller, nil
}

func AssignVehicle(dispatchID, vehicleID, operator string) (*models.Dispatch, error) {
	dispatch, ok := store.GlobalStore.GetDispatch(dispatchID)
	if !ok {
		return nil, ErrDispatchNotFound
	}

	if dispatch.Status != models.DispatchStatusQueued {
		return nil, ErrInvalidStatusTransition
	}

	_, ok := store.GlobalStore.AssignVehicleToDispatch(vehicleID, dispatchID)
	if !ok {
		vehicle, exists := store.GlobalStore.GetVehicle(vehicleID)
		if !exists {
			return nil, ErrVehicleNotFound
		}
		if vehicle.Status != models.VehicleStatusStandby {
			return nil, ErrVehicleOccupied
		}
		return nil, ErrVehicleOccupied
	}

	now := time.Now()

	dispatch.VehicleID = vehicleID
	dispatch.Status = models.DispatchStatusDispatched
	dispatch.DispatchTime = &now
	dispatch.ChangeLog = append(dispatch.ChangeLog, models.ChangeLogEntry{
		Timestamp: now,
		Action:    "assigned",
		Operator:  operator,
		OldValue:  "",
		NewValue:  vehicleID,
		Reason:    "派车",
	})

	store.GlobalStore.UpdateDispatch(dispatch)
	store.GlobalStore.RemoveFromQueue(dispatchID)

	caller, _ := store.GlobalStore.GetCaller(dispatch.CallerID)
	if caller != nil {
		caller.DispatchIDs = append(caller.DispatchIDs, dispatchID)
		store.GlobalStore.UpdateCaller(caller)
	}

	return dispatch, nil
}

func DispatchArriveOnScene(dispatchID, operator string) (*models.Dispatch, error) {
	dispatch, ok := store.GlobalStore.GetDispatch(dispatchID)
	if !ok {
		return nil, ErrDispatchNotFound
	}

	if dispatch.Status != models.DispatchStatusDispatched {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	dispatch.Status = models.DispatchStatusOnScene
	dispatch.ArrivalTime = &now
	dispatch.ChangeLog = append(dispatch.ChangeLog, models.ChangeLogEntry{
		Timestamp: now,
		Action:    "arrive_on_scene",
		Operator:  operator,
		OldValue:  string(models.DispatchStatusDispatched),
		NewValue:  string(models.DispatchStatusOnScene),
	})

	vehicle, ok := store.GlobalStore.GetVehicle(dispatch.VehicleID)
	if ok {
		vehicle.Status = models.VehicleStatusOnScene
		vehicle.LastStatusUpdate = now
		store.GlobalStore.UpdateVehicle(vehicle)
	}

	store.GlobalStore.UpdateDispatch(dispatch)
	return dispatch, nil
}

func DispatchStartTransfer(dispatchID, hospitalID, operator string) (*models.Dispatch, error) {
	dispatch, ok := store.GlobalStore.GetDispatch(dispatchID)
	if !ok {
		return nil, ErrDispatchNotFound
	}

	if dispatch.Status != models.DispatchStatusOnScene {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	dispatch.Status = models.DispatchStatusInTransit
	dispatch.HospitalID = hospitalID
	dispatch.DepartSceneTime = &now
	dispatch.ChangeLog = append(dispatch.ChangeLog, models.ChangeLogEntry{
		Timestamp: now,
		Action:    "start_transfer",
		Operator:  operator,
		OldValue:  string(models.DispatchStatusOnScene),
		NewValue:  string(models.DispatchStatusInTransit),
		Reason:    fmt.Sprintf("转送目标医院: %s", hospitalID),
	})

	vehicle, ok := store.GlobalStore.GetVehicle(dispatch.VehicleID)
	if ok {
		vehicle.Status = models.VehicleStatusInTransit
		vehicle.LastStatusUpdate = now
		store.GlobalStore.UpdateVehicle(vehicle)
	}

	store.GlobalStore.UpdateDispatch(dispatch)
	return dispatch, nil
}

func DispatchReturn(dispatchID, operator string, isEmptyRun bool, emptyRunReason string) (*models.Dispatch, error) {
	dispatch, ok := store.GlobalStore.GetDispatch(dispatchID)
	if !ok {
		return nil, ErrDispatchNotFound
	}

	if dispatch.Status != models.DispatchStatusInTransit && dispatch.Status != models.DispatchStatusOnScene {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	dispatch.Status = models.DispatchStatusCompleted
	dispatch.ReturnTime = &now
	dispatch.CompleteTime = &now
	dispatch.IsEmptyRun = isEmptyRun
	if isEmptyRun {
		dispatch.EmptyRunReason = emptyRunReason
	}
	dispatch.ChangeLog = append(dispatch.ChangeLog, models.ChangeLogEntry{
		Timestamp: now,
		Action:    "completed",
		Operator:  operator,
		OldValue:  string(dispatch.Status),
		NewValue:  string(models.DispatchStatusCompleted),
		Reason:    "任务完成，车辆返站",
	})

	vehicle, ok := store.GlobalStore.GetVehicle(dispatch.VehicleID)
	if ok {
		vehicle.Status = models.VehicleStatusReturning
		vehicle.LastStatusUpdate = now
		if dispatch.DispatchTime != nil {
			drivingDuration := now.Sub(*dispatch.DispatchTime).Seconds()
			vehicle.TotalDrivingSeconds += int64(drivingDuration)
		}
		vehicle.CurrentDispatchID = ""
		store.GlobalStore.UpdateVehicle(vehicle)
	}

	store.GlobalStore.UpdateDispatch(dispatch)
	return dispatch, nil
}

func VehicleToStandby(vehicleID, operator string) (*models.Vehicle, error) {
	vehicle, ok := store.GlobalStore.GetVehicle(vehicleID)
	if !ok {
		return nil, ErrVehicleNotFound
	}

	if vehicle.Status != models.VehicleStatusReturning && vehicle.Status != models.VehicleStatusMaintenance {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	vehicle.Status = models.VehicleStatusStandby
	vehicle.MaintenanceType = ""
	vehicle.LastStatusUpdate = now
	store.GlobalStore.UpdateVehicle(vehicle)

	return vehicle, nil
}

func VehicleToMaintenance(vehicleID, operator string, mType models.MaintenanceType) (*models.Vehicle, error) {
	vehicle, ok := store.GlobalStore.GetVehicle(vehicleID)
	if !ok {
		return nil, ErrVehicleNotFound
	}

	if vehicle.Status != models.VehicleStatusReturning && vehicle.Status != models.VehicleStatusStandby {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	vehicle.Status = models.VehicleStatusMaintenance
	vehicle.MaintenanceType = mType
	vehicle.LastStatusUpdate = now
	store.GlobalStore.UpdateVehicle(vehicle)

	return vehicle, nil
}

func ReassignVehicle(dispatchID, newVehicleID, operator, reason string) (*models.Dispatch, error) {
	dispatch, ok := store.GlobalStore.GetDispatch(dispatchID)
	if !ok {
		return nil, ErrDispatchNotFound
	}

	if dispatch.Status == models.DispatchStatusCompleted ||
		dispatch.Status == models.DispatchStatusCancelled ||
		dispatch.Status == models.DispatchStatusReassigned {
		return nil, ErrInvalidStatusTransition
	}

	now := time.Now()
	oldVehicleID := dispatch.VehicleID

	if oldVehicleID == newVehicleID {
		return dispatch, nil
	}

	_, ok = store.GlobalStore.AssignVehicleToDispatch(newVehicleID, dispatchID)
	if !ok {
		return nil, ErrVehicleOccupied
	}

	if oldVehicleID != "" {
		store.GlobalStore.ReleaseVehicleFromDispatch(oldVehicleID)
	}

	originalDispatchID := dispatch.OriginalDispatchID
	if originalDispatchID == "" {
		originalDispatchID = dispatchID
	}

	dispatch.VehicleID = newVehicleID
	dispatch.ReassignedFromID = oldVehicleID
	dispatch.OriginalDispatchID = originalDispatchID
	dispatch.ChangeLog = append(dispatch.ChangeLog, models.ChangeLogEntry{
		Timestamp: now,
		Action:    "reassign",
		Operator:  operator,
		OldValue:  oldVehicleID,
		NewValue:  newVehicleID,
		Reason:    reason,
	})

	store.GlobalStore.UpdateDispatch(dispatch)
	return dispatch, nil
}
