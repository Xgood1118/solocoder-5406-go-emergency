package handler

import (
	"net/http"
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/service"
	"emergency-dispatch/internal/store"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ListVehicles(c *gin.Context) {
	status := c.Query("status")
	var vehicles []*models.Vehicle

	if status != "" {
		vehicles = store.GlobalStore.ListVehiclesByStatus(models.VehicleStatus(status))
	} else {
		vehicles = store.GlobalStore.ListVehicles()
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    vehicles,
	})
}

func GetVehicle(c *gin.Context) {
	id := c.Param("id")
	vehicle, ok := store.GlobalStore.GetVehicle(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "车辆不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    vehicle,
	})
}

type CreateVehicleRequest struct {
	PlateNumber string             `json:"plate_number"`
	Type        models.VehicleType `json:"type"`
	Crew        []string           `json:"crew"`
	Longitude   float64            `json:"longitude"`
	Latitude    float64            `json:"latitude"`
}

func CreateVehicle(c *gin.Context) {
	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	now := time.Now()
	vehicleID := fmt.Sprintf("V%d", now.UnixNano())
	vehicle := &models.Vehicle{
		ID:               vehicleID,
		PlateNumber:      req.PlateNumber,
		Type:             req.Type,
		Crew:             req.Crew,
		Longitude:        req.Longitude,
		Latitude:         req.Latitude,
		Status:           models.VehicleStatusStandby,
		LastStatusUpdate: now,
		OnlineSince:      now,
	}

	store.GlobalStore.AddVehicle(vehicle)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    vehicle,
	})
}

type UpdateVehicleGPSRequest struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func UpdateVehicleGPS(c *gin.Context) {
	id := c.Param("id")
	var req UpdateVehicleGPSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	vehicle, ok := store.GlobalStore.GetVehicle(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "车辆不存在",
		})
		return
	}

	vehicle.Longitude = req.Longitude
	vehicle.Latitude = req.Latitude
	store.GlobalStore.UpdateVehicle(vehicle)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    vehicle,
	})
}

type VehicleToStandbyRequest struct {
	Operator string `json:"operator"`
}

func VehicleToStandbyHandler(c *gin.Context) {
	id := c.Param("id")
	var req VehicleToStandbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	vehicle, err := service.VehicleToStandby(id, req.Operator)
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
		Data:    vehicle,
	})
}

type VehicleToMaintenanceRequest struct {
	Operator      string                `json:"operator"`
	MaintenanceType models.MaintenanceType `json:"maintenance_type"`
}

func VehicleToMaintenanceHandler(c *gin.Context) {
	id := c.Param("id")
	var req VehicleToMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	vehicle, err := service.VehicleToMaintenance(id, req.Operator, req.MaintenanceType)
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
		Data:    vehicle,
	})
}
