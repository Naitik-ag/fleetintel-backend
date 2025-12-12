package models

import "time"

type DeliveryAgent struct {
	ID          string    `gorm:"primaryKey" json:"id"`                      // AGENT001, AGENT002, etc.
	Name        string    `gorm:"not null" json:"name"`                      // Full name
	Phone       string    `gorm:"unique;not null" json:"phone"`              // Contact number (unique)
	Email       string    `gorm:"unique" json:"email"`                       // Email address
	VehicleType string    `gorm:"type:varchar(20)" json:"vehicle_type"`      // bike, scooter, car, truck
	Status      string    `gorm:"type:varchar(20);default:'offline'" json:"status"` // available, busy, offline
	IsActive    bool      `gorm:"default:true" json:"is_active"`             // Account active/inactive
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`          // When agent registered
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`          // Last profile update
}

func (DeliveryAgent) TableName() string {
	return "delivery_agents"
}

type AgentRequest struct {
	ID          string `json:"id" validate:"required"`                    // Agent ID (required)
	Name        string `json:"name" validate:"required"`                  // Name (required)
	Phone       string `json:"phone" validate:"required"`                 // Phone (required)
	Email       string `json:"email" validate:"email"`                    // Email (optional, must be valid format)
	VehicleType string `json:"vehicle_type" validate:"required"`          // Vehicle type (required)
}

type AgentResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	VehicleType string    `json:"vehicle_type"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AgentStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=available busy offline"` // Only these values allowed
}

type AgentStats struct {
	AgentID          string  `json:"agent_id"`
	TotalDeliveries  int     `json:"total_deliveries"`
	TotalDistance    float64 `json:"total_distance_km"`
	AverageRating    float64 `json:"average_rating"`
	TotalEarnings    float64 `json:"total_earnings"`
	ActiveSince      string  `json:"active_since"`
}
