package models

import "time"

type LocationRequest struct {
	AgentID   string  `json:"agent_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp string  `json:"timestamp"`
}

type Location struct {
	ID        uint      `gorm:"primaryKey"`
	AgentID   string    `gorm:"index;not null"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null"`
	Longitude float64   `gorm:"type:decimal(11,8);not null"`
	Speed     float64   `gorm:"type:decimal(6,2)"`
	Heading   float64   `gorm:"type:decimal(5,2)"`
	Accuracy  float64   `gorm:"type:decimal(6,2)"`
	Status    string    `gorm:"type:varchar(20);default:'unknown'"`
	Timestamp time.Time `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Location) TableName() string {
	return "locations"
}

type LocationResponse struct {
	ID        uint      `json:"id"`
	AgentID   string    `json:"agent_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Accuracy  float64   `json:"accuracy"`
	Altitude  float64   `json:"altitude"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}