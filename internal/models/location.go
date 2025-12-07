package models

import "time"

type LocationRequest struct {
	VehicleID string  `json:"vehicle_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Speed     float64 `json:"speed" validate:"min=0"`
	Heading   float64 `json:"heading" validate:"min=0,max=360"`
	Accuracy  float64 `json:"accuracy" validate:"min=0"`
	Timestamp string  `json:"timestamp" validate:"required"`
}

type Location struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	VehicleID string    `json:"vehicle_id" gorm:"index;not null"`
	Latitude  float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Accuracy  float64   `json:"accuracy"`
	Timestamp time.Time `json:"timestamp" gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type LocationResponse struct {
	VehicleID   string    `json:"vehicle_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Speed       float64   `json:"speed"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}
