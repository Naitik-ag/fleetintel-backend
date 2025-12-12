package handlers

import (
	"log"
	"strings"

	"github.com/Naitik-ag/fleetintel-backend/internal/models"
	"github.com/Naitik-ag/fleetintel-backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type AgentHandler struct {
	agentRepo *repository.AgentRepository
}

func NewAgentHandler(agentRepo *repository.AgentRepository) *AgentHandler {
	return &AgentHandler{
		agentRepo: agentRepo,
	}
}

func (h *AgentHandler) RegisterAgent(c *fiber.Ctx) error {
	var req models.AgentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.ID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "agent_id is required",
		})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if req.Phone == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "phone is required",
		})
	}

	if req.VehicleType == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "vehicle_type is required (bike, scooter, car, truck)",
		})
	}

	validVehicles := map[string]bool{
		"bike":    true,
		"scooter": true,
		"car":     true,
		"truck":   true,
	}

	if !validVehicles[strings.ToLower(req.VehicleType)] {
		return c.Status(400).JSON(fiber.Map{
			"error": "vehicle_type must be: bike, scooter, car, or truck",
		})
	}

	agent := models.DeliveryAgent{
		ID:          req.ID,
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		VehicleType: strings.ToLower(req.VehicleType),
		Status:      "offline",
		IsActive:    true,
	}

	err := h.agentRepo.Create(&agent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to register agent",
			"details": err.Error(),
		})
	}

	response := models.AgentResponse{
		ID:          agent.ID,
		Name:        agent.Name,
		Phone:       agent.Phone,
		Email:       agent.Email,
		VehicleType: agent.VehicleType,
		Status:      agent.Status,
		IsActive:    agent.IsActive,
		CreatedAt:   agent.CreatedAt,
		UpdatedAt:   agent.UpdatedAt,
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Agent registered successfully",
		"data":    response,
	})
}

func (h *AgentHandler) GetAgent(c *fiber.Ctx) error {
	agentID := c.Params("id")

	if agentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "agent_id parameter is required",
		})
	}

	agent, err := h.agentRepo.FindByID(agentID)
	if err != nil {
		log.Printf("Failed to fetch agent: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch agent",
			"details": err.Error(),
		})
	}

	if agent == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Agent not found",
		})
	}

	response := models.AgentResponse{
		ID:          agent.ID,
		Name:        agent.Name,
		Phone:       agent.Phone,
		Email:       agent.Email,
		VehicleType: agent.VehicleType,
		Status:      agent.Status,
		IsActive:    agent.IsActive,
		CreatedAt:   agent.CreatedAt,
		UpdatedAt:   agent.UpdatedAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

func (h *AgentHandler) ListAgents(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	page := c.QueryInt("page", 1)
	status := c.Query("status", "")

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	agents, totalCount, err := h.agentRepo.FindAll(limit, offset, status)
	if err != nil {
		log.Printf("Failed to fetch agents: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch agents",
			"details": err.Error(),
		})
	}

	var responses []models.AgentResponse
	for _, agent := range agents {
		responses = append(responses, models.AgentResponse{
			ID:          agent.ID,
			Name:        agent.Name,
			Phone:       agent.Phone,
			Email:       agent.Email,
			VehicleType: agent.VehicleType,
			Status:      agent.Status,
			IsActive:    agent.IsActive,
			CreatedAt:   agent.CreatedAt,
			UpdatedAt:   agent.UpdatedAt,
		})
	}

	totalPages := int(totalCount) / limit
	if int(totalCount)%limit != 0 {
		totalPages++
	}

	log.Printf("Retrieved %d agents (page %d of %d)", len(agents), page, totalPages)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
		"pagination": fiber.Map{
			"total":        totalCount,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	})
}

func (h *AgentHandler) UpdateAgent(c *fiber.Ctx) error {
	agentID := c.Params("id")
	var req models.AgentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.VehicleType != "" {
		validVehicles := map[string]bool{
			"bike": true, "scooter": true, "car": true, "truck": true,
		}
		vehicleType := strings.ToLower(req.VehicleType)

		if !validVehicles[vehicleType] {
			return c.Status(400).JSON(fiber.Map{
				"error": "vehicle_type must be: bike, scooter, car, or truck",
			})
		}
		updates["vehicle_type"] = vehicleType
	}

	if len(updates) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	err := h.agentRepo.Update(agentID, updates)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update agent",
			"details": err.Error(),
		})
	}

	agent, _ := h.agentRepo.FindByID(agentID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Agent updated successfully",
		"data": models.AgentResponse{
			ID:          agent.ID,
			Name:        agent.Name,
			Phone:       agent.Phone,
			Email:       agent.Email,
			VehicleType: agent.VehicleType,
			Status:      agent.Status,
			IsActive:    agent.IsActive,
			CreatedAt:   agent.CreatedAt,
			UpdatedAt:   agent.UpdatedAt,
		},
	})
}

func (h *AgentHandler) UpdateAgentStatus(c *fiber.Ctx) error {
	agentID := c.Params("id")
	var req models.AgentStatusUpdate

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.Status == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "status is required (available, busy, offline)",
		})
	}

	err := h.agentRepo.UpdateStatus(agentID, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update status",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Agent status updated successfully",
		"data": fiber.Map{
			"agent_id": agentID,
			"status":   req.Status,
		},
	})
}

func (h *AgentHandler) DeleteAgent(c *fiber.Ctx) error {
	agentID := c.Params("id")

	err := h.agentRepo.SoftDelete(agentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete agent",
			"details": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Agent deleted successfully",
	})
}

func (h *AgentHandler) GetAgentStats(c *fiber.Ctx) error {
	agentID := c.Params("id")

	agent, err := h.agentRepo.FindByID(agentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch agent",
		})
	}

	if agent == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Agent not found",
		})
	}

	stats := models.AgentStats{
		AgentID:         agentID,
		TotalDeliveries: 0,
		TotalDistance:   0.0,
		AverageRating:   0.0,
		TotalEarnings:   0.0,
		ActiveSince:     agent.CreatedAt.Format("2006-01-02"),
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}
