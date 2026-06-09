package main

import (
	"log"
	"os"

	"emergency-dispatch/internal/handler"
	"emergency-dispatch/internal/seed"

	"github.com/gin-gonic/gin"
)

func main() {
	seed.SeedData()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		vehicles := api.Group("/vehicles")
		{
			vehicles.GET("", handler.ListVehicles)
			vehicles.GET("/:id", handler.GetVehicle)
			vehicles.POST("", handler.CreateVehicle)
			vehicles.PUT("/:id/gps", handler.UpdateVehicleGPS)
			vehicles.POST("/:id/standby", handler.VehicleToStandbyHandler)
			vehicles.POST("/:id/maintenance", handler.VehicleToMaintenanceHandler)
		}

		hospitals := api.Group("/hospitals")
		{
			hospitals.GET("", handler.ListHospitals)
			hospitals.GET("/:id", handler.GetHospital)
			hospitals.POST("", handler.CreateHospital)
			hospitals.PUT("/:id/beds", handler.UpdateHospitalBeds)
		}

		dispatches := api.Group("/dispatches")
		{
			dispatches.GET("", handler.ListDispatches)
			dispatches.GET("/queue", handler.GetDispatchQueue)
			dispatches.GET("/:id", handler.GetDispatch)
			dispatches.POST("", handler.CreateDispatchHandler)
			dispatches.POST("/recommend", handler.RecommendVehicles)
			dispatches.POST("/:id/assign", handler.AssignVehicleHandler)
			dispatches.POST("/:id/arrive", handler.DispatchArriveOnSceneHandler)
			dispatches.POST("/:id/transfer", handler.DispatchStartTransferHandler)
			dispatches.POST("/:id/return", handler.DispatchReturnHandler)
			dispatches.POST("/:id/reassign", handler.ReassignVehicleHandler)
		}

		pois := api.Group("/pois")
		{
			pois.GET("", handler.ListPOIs)
			pois.GET("/match", handler.MatchPOIHandler)
			pois.GET("/:id", handler.GetPOI)
			pois.POST("", handler.CreatePOI)
		}

		callers := api.Group("/callers")
		{
			callers.GET("", handler.ListCallers)
			callers.GET("/search", handler.GetCallerByPhone)
			callers.GET("/:id", handler.GetCaller)
			callers.PUT("/:id/vip", handler.UpdateCallerVIP)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/monthly", handler.GetMonthlyStatsHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
