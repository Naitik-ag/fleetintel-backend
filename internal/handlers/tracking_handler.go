package handlers

import (
	"log"
	"time"

	"github.com/Naitik-ag/fleetintel-backend/internal/models"
	"github.com/Naitik-ag/fleetintel-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type TrackingHandler struct {
	locationRepo *repository.LocationRepository
}

func NewTrackingHandler(locationRepo *repository.LocationRepository) *TrackingHandler {
	return &TrackingHandler{
		locationRepo: locationRepo,
	}
}

func calculateStatus(speed float64) string {
	switch {
	case speed == 0:
		return "stopped"
	case speed > 0 && speed < 5:
		return "idle"
	case speed >= 5:
		return "moving"
	default:
		return "unknown"
	}
}

func (h *TrackingHandler) UpdateLocation(c *fiber.Ctx) error {
	var req models.LocationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.AgentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "agent_id is required",
		})
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "latitude and longitude are required",
		})
	}

	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	status := calculateStatus(req.Speed)

	location := models.Location{
		AgentID:   req.AgentID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Accuracy:  req.Accuracy,
		Status:    status,
		Timestamp: timestamp,
	}

	if err := h.locationRepo.Create(&location); err != nil {
		log.Printf("Failed to save location: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to save location",
			"details": err.Error(),
		})
	}

	log.Printf("Location saved: Agent=%s, Status=%s, Lat=%.6f, Lng=%.6f, Speed=%.2f km/h, ID=%d",
		req.AgentID, status, req.Latitude, req.Longitude, req.Speed, location.ID) // 🔄 CHANGED log

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Location updated successfully",
		"data": fiber.Map{
			"id":        location.ID,
			"agent_id":  req.AgentID,
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
			"status":    status,
			"timestamp": location.Timestamp,
		},
	})
}

func (h *TrackingHandler) GetLiveLocation(c *fiber.Ctx) error {
	agentID := c.Params("id")

	if agentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "agent_id parameter is required",
		})
	}

	location, err := h.locationRepo.FindLatestByAgentID(agentID)
	if err != nil {
		log.Printf("Failed to fetch live location: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch location",
			"details": err.Error(),
		})
	}

	if location == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "No location found for this agent",
		})
	}

	log.Printf("Fetched live location for agent: %s (Status: %s)", agentID, location.Status)

	response := models.LocationResponse{
		ID:        location.ID,
		AgentID:   location.AgentID,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		Speed:     location.Speed,
		Heading:   location.Heading,
		Accuracy:  location.Accuracy,
		Status:    location.Status,
		Timestamp: location.Timestamp,
		CreatedAt: location.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

func (h *TrackingHandler) GetLocationHistory(c *fiber.Ctx) error {
	agentID := c.Params("id")

	if agentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "agent_id parameter is required",
		})
	}

	from := c.Query("from")
	to := c.Query("to")
	limit := c.QueryInt("limit", 100)

	if limit > 1000 {
		limit = 1000
	}

	log.Printf("Fetching history for agent: %s (from=%s, to=%s, limit=%d)",
		agentID, from, to, limit)

	var locations []models.Location
	var err error

	if from != "" && to != "" {
		startTime, err1 := time.Parse(time.RFC3339, from)
		endTime, err2 := time.Parse(time.RFC3339, to)

		if err1 != nil || err2 != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid time format. Use RFC3339 format: 2024-12-07T10:30:00Z",
			})
		}

		locations, err = h.locationRepo.FindByAgentIDAndTimeRange(agentID, startTime, endTime)
	} else {
		locations, err = h.locationRepo.FindByAgentID(agentID, limit)
	}

	if err != nil {
		log.Printf("Failed to fetch location history: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch location history",
			"details": err.Error(),
		})
	}

	if len(locations) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "No location history found for this agent",
			"data":    []models.Location{},
		})
	}

	var responses []models.LocationResponse
	for _, loc := range locations {
		responses = append(responses, models.LocationResponse{
			ID:        loc.ID,
			AgentID:   loc.AgentID,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Speed:     loc.Speed,
			Heading:   loc.Heading,
			Accuracy:  loc.Accuracy,
			Status:    loc.Status,
			Timestamp: loc.Timestamp,
			CreatedAt: loc.CreatedAt,
		})
	}

	log.Printf("Retrieved %d locations for agent: %s", len(locations), agentID)
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(responses),
		"data":    responses,
	})
}
