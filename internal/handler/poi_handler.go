package handler

import (
	"fmt"
	"net/http"
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/service"
	"emergency-dispatch/internal/store"

	"github.com/gin-gonic/gin"
)

func ListPOIs(c *gin.Context) {
	pois := store.GlobalStore.ListPOIs()
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    pois,
	})
}

func GetPOI(c *gin.Context) {
	id := c.Param("id")
	poi, ok := store.GlobalStore.GetPOI(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "POI 不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    poi,
	})
}

type CreatePOIRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Type      string  `json:"type"`
}

func CreatePOI(c *gin.Context) {
	var req CreatePOIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	now := time.Now()
	poiID := fmt.Sprintf("P%d", now.UnixNano())
	poi := &models.POI{
		ID:        poiID,
		Name:      req.Name,
		Address:   req.Address,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		Type:      req.Type,
	}

	store.GlobalStore.AddPOI(poi)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    poi,
	})
}

func MatchPOIHandler(c *gin.Context) {
	query := c.Query("q")
	limit := 10
	if l := c.Query("limit"); l != "" {
		if n, err := fmt.Sscanf(l, "%d", &limit); err == nil && n == 1 {
			if limit <= 0 {
				limit = 10
			}
			if limit > 50 {
				limit = 50
			}
		}
	}

	results := service.MatchPOIs(query, limit)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    results,
	})
}
