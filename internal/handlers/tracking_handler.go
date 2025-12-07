package handlers

import (
	"log"
	"time"

	"github.com/Naitik-ag/fleetintel-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type TrackingHandler struct {
}

func NewTrackingHandler() *TrackingHandler {
	return &TrackingHandler{}
}

func (h *TrackingHandler) UpdateLocation(c *fiber.Ctx) error {
	var req models.LocationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.VehicleID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "vehicle_id is required",
		})
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "latitude and longitude are required",
		})
	}

	// TODO : Call service layer to process and save
	// err := h.trackingService.SaveLocation(req)

	log.Printf("Location received: Vehicle=%s, Lat=%.6f, Lng=%.6f, Speed=%.2f km/h",
		req.VehicleID, req.Latitude, req.Longitude, req.Speed)

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Location updated successfully",
		"data": fiber.Map{
			"vehicle_id": req.VehicleID,
			"latitude":   req.Latitude,
			"longitude":  req.Longitude,
		},
	})
}

func (h *TrackingHandler) GetLiveLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("id")

	if vehicleID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "vehicle_id parameter is required",
		})
	}

	// TODO: Fetch from database
	// location, err := h.repository.GetLatestLocation(vehicleID)

	log.Printf("Fetching live location for vehicle: %s", vehicleID)

	// Mock response using our model
	mockResponse := models.LocationResponse{
		VehicleID:   vehicleID,
		Latitude:    28.6139,
		Longitude:   77.2090,
		Speed:       45.5,
		Status:      "moving",
		LastUpdated: time.Now(),
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    mockResponse,
	})
}

func (h *TrackingHandler) GetLocationHistory(c *fiber.Ctx) error {
	vehicleID := c.Params("id")

	from := c.Query("from")
	to := c.Query("to")
	limit := c.QueryInt("limit", 100)

	log.Printf("Fetching history for vehicle: %s (from=%s, to=%s, limit=%d)",
		vehicleID, from, to, limit)

	// TODO : Fetch from database with filters
	// locations, err := h.repository.GetHistory(vehicleID, from, to, limit)

	// Mock response
	mockLocations := []models.LocationResponse{
		{
			VehicleID:   vehicleID,
			Latitude:    28.6139,
			Longitude:   77.2090,
			Speed:       45.5,
			Status:      "moving",
			LastUpdated: time.Now().Add(-5 * time.Minute),
		},
		{
			VehicleID:   vehicleID,
			Latitude:    28.6200,
			Longitude:   77.2100,
			Speed:       50.0,
			Status:      "moving",
			LastUpdated: time.Now().Add(-10 * time.Minute),
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(mockLocations),
		"data":    mockLocations,
	})
}
