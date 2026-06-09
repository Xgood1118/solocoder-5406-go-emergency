package handler

import (
	"fmt"
	"net/http"
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"

	"github.com/gin-gonic/gin"
)

func ListHospitals(c *gin.Context) {
	hospitals := store.GlobalStore.ListHospitals()
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    hospitals,
	})
}

func GetHospital(c *gin.Context) {
	id := c.Param("id")
	hospital, ok := store.GlobalStore.GetHospital(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "医院不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    hospital,
	})
}

type CreateHospitalRequest struct {
	Name      string                        `json:"name"`
	Address   string                        `json:"address"`
	Longitude float64                       `json:"longitude"`
	Latitude  float64                       `json:"latitude"`
	Beds      map[models.HospitalDepartment]int `json:"beds"`
}

func CreateHospital(c *gin.Context) {
	var req CreateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	now := time.Now()
	hospitalID := fmt.Sprintf("H%d", now.UnixNano())
	hospital := &models.Hospital{
		ID:        hospitalID,
		Name:      req.Name,
		Address:   req.Address,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		Beds:      req.Beds,
	}

	if hospital.Beds == nil {
		hospital.Beds = make(map[models.HospitalDepartment]int)
	}

	store.GlobalStore.AddHospital(hospital)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    hospital,
	})
}

type UpdateHospitalBedsRequest struct {
	Beds map[models.HospitalDepartment]int `json:"beds"`
}

func UpdateHospitalBeds(c *gin.Context) {
	id := c.Param("id")
	var req UpdateHospitalBedsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	hospital, ok := store.GlobalStore.GetHospital(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "医院不存在",
		})
		return
	}

	for dept, count := range req.Beds {
		hospital.Beds[dept] = count
	}

	store.GlobalStore.UpdateHospital(hospital)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    hospital,
	})
}
