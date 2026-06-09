package handler

import (
	"net/http"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/service"
	"emergency-dispatch/internal/store"

	"github.com/gin-gonic/gin"
)

func ListDispatches(c *gin.Context) {
	status := c.Query("status")
	var dispatches []*models.Dispatch

	if status != "" {
		dispatches = store.GlobalStore.ListDispatchesByStatus(models.DispatchStatus(status))
	} else {
		dispatches = store.GlobalStore.ListDispatches()
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatches,
	})
}

func GetDispatch(c *gin.Context) {
	id := c.Param("id")
	dispatch, ok := store.GlobalStore.GetDispatch(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "调度记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type CreateDispatchRequest struct {
	CallerName         string             `json:"caller_name"`
	CallerPhone        string             `json:"caller_phone"`
	IsVIP              bool               `json:"is_vip"`
	QueueJumpReason    string             `json:"queue_jump_reason"`
	IllnessType        models.IllnessType `json:"illness_type"`
	IllnessDescription string             `json:"illness_description"`
	Address            string             `json:"address"`
	Longitude          float64            `json:"longitude"`
	Latitude           float64            `json:"latitude"`
	Operator           string             `json:"operator"`
}

func CreateDispatchHandler(c *gin.Context) {
	var req CreateDispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.CreateDispatch(&service.CreateDispatchRequest{
		CallerName:         req.CallerName,
		CallerPhone:        req.CallerPhone,
		IsVIP:              req.IsVIP,
		QueueJumpReason:    req.QueueJumpReason,
		IllnessType:        req.IllnessType,
		IllnessDescription: req.IllnessDescription,
		Address:            req.Address,
		Longitude:          req.Longitude,
		Latitude:           req.Latitude,
		Operator:           req.Operator,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type AssignVehicleRequest struct {
	VehicleID string `json:"vehicle_id"`
	Operator  string `json:"operator"`
}

func AssignVehicleHandler(c *gin.Context) {
	id := c.Param("id")
	var req AssignVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.AssignVehicle(id, req.VehicleID, req.Operator)
	if err != nil {
		code := 400
		if err == service.ErrVehicleOccupied {
			code = 409
		}
		c.JSON(code, models.Response{
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type ArriveOnSceneRequest struct {
	Operator string `json:"operator"`
}

func DispatchArriveOnSceneHandler(c *gin.Context) {
	id := c.Param("id")
	var req ArriveOnSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.DispatchArriveOnScene(id, req.Operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type StartTransferRequest struct {
	HospitalID string `json:"hospital_id"`
	Operator   string `json:"operator"`
}

func DispatchStartTransferHandler(c *gin.Context) {
	id := c.Param("id")
	var req StartTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.DispatchStartTransfer(id, req.HospitalID, req.Operator)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type DispatchReturnRequest struct {
	Operator       string `json:"operator"`
	IsEmptyRun     bool   `json:"is_empty_run"`
	EmptyRunReason string `json:"empty_run_reason"`
}

func DispatchReturnHandler(c *gin.Context) {
	id := c.Param("id")
	var req DispatchReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.DispatchReturn(id, req.Operator, req.IsEmptyRun, req.EmptyRunReason)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type ReassignVehicleRequest struct {
	NewVehicleID string `json:"new_vehicle_id"`
	Operator     string `json:"operator"`
	Reason       string `json:"reason"`
}

func ReassignVehicleHandler(c *gin.Context) {
	id := c.Param("id")
	var req ReassignVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	dispatch, err := service.ReassignVehicle(id, req.NewVehicleID, req.Operator, req.Reason)
	if err != nil {
		code := 400
		if err == service.ErrVehicleOccupied {
			code = 409
		}
		c.JSON(code, models.Response{
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatch,
	})
}

type RecommendVehiclesRequest struct {
	IllnessType models.IllnessType `json:"illness_type"`
	Longitude   float64            `json:"longitude"`
	Latitude    float64            `json:"latitude"`
	Limit       int                `json:"limit"`
}

func RecommendVehicles(c *gin.Context) {
	var req RecommendVehiclesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 5
	}

	candidates := service.FindNearestVehicles(req.Latitude, req.Longitude, req.IllnessType, req.Limit)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    candidates,
	})
}

func GetDispatchQueue(c *gin.Context) {
	queue := store.GlobalStore.GetQueue()
	var dispatches []*models.Dispatch
	for _, id := range queue {
		d, ok := store.GlobalStore.GetDispatch(id)
		if ok {
			dispatches = append(dispatches, d)
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    dispatches,
	})
}
