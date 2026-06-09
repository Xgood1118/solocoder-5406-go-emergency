package handler

import (
	"net/http"
	"strconv"
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/service"

	"github.com/gin-gonic/gin"
)

func GetMonthlyStatsHandler(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			if m >= 1 && m <= 12 {
				month = m
			}
		}
	}

	stats := service.GetMonthlyStats(year, month)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    stats,
	})
}
