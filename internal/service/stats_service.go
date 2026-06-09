package service

import (
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"
)

func GetMonthlyStats(year, month int) *models.MonthlyStats {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)

	dispatches := store.GlobalStore.ListDispatches()
	var monthDispatches []*models.Dispatch

	for _, d := range dispatches {
		if d.CallTime.After(start) && d.CallTime.Before(end) {
			monthDispatches = append(monthDispatches, d)
		}
	}

	total := len(monthDispatches)
	emptyRunCount := 0
	var totalResponseSeconds float64
	responseCount := 0
	hospitalStats := make(map[string]int)

	for _, d := range monthDispatches {
		if d.IsEmptyRun {
			emptyRunCount++
		}
		if d.ArrivalTime != nil {
			respTime := d.ArrivalTime.Sub(d.CallTime).Seconds()
			totalResponseSeconds += respTime
			responseCount++
		}
		if d.HospitalID != "" && d.Status == models.DispatchStatusCompleted {
			hospitalStats[d.HospitalID]++
		}
	}

	avgResponseTime := 0.0
	if responseCount > 0 {
		avgResponseTime = totalResponseSeconds / float64(responseCount)
	}

	emptyRunRate := 0.0
	if total > 0 {
		emptyRunRate = float64(emptyRunCount) / float64(total)
	}

	vehicles := store.GlobalStore.ListVehicles()
	utilization := make(map[string]float64)
	now := time.Now()

	for _, v := range vehicles {
		onlineDuration := now.Sub(v.OnlineSince).Seconds()
		if onlineDuration > 0 {
			util := float64(v.TotalDrivingSeconds) / onlineDuration
			if util > 1.0 {
				util = 1.0
			}
			utilization[v.ID] = util
		} else {
			utilization[v.ID] = 0.0
		}
	}

	return &models.MonthlyStats{
		Month:               start.Format("2006-01"),
		AvgResponseTimeSec:  avgResponseTime,
		TotalDispatches:     total,
		EmptyRunCount:       emptyRunCount,
		EmptyRunRate:        emptyRunRate,
		HospitalReceiveStats: hospitalStats,
		VehicleUtilization:  utilization,
	}
}
