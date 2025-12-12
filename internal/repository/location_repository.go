package repository

import (
	"time"

	"gorm.io/gorm"
	"github.com/Naitik-ag/fleetintel-backend/internal/models"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) Create(location *models.Location) error {
	result := r.db.Create(location)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *LocationRepository) FindByAgentID(agentID string, limit int) ([]models.Location, error) {
	var locations []models.Location
	
	result := r.db.Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&locations)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return locations, nil
}

func (r *LocationRepository) FindLatestByAgentID(agentID string) (*models.Location, error) {
	var location models.Location
	
	result := r.db.Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		First(&location)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &location, nil
}


func (r *LocationRepository) FindByAgentIDAndTimeRange(agentID string, startTime, endTime time.Time) ([]models.Location, error) {
	var locations []models.Location
	
	result := r.db.Where("agent_id = ?", agentID).
		Where("timestamp >= ?", startTime).
		Where("timestamp <= ?", endTime).
		Order("timestamp ASC").
		Find(&locations)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return locations, nil
}

func (r *LocationRepository) Count(agentID string) (int64, error) {
	var count int64
	
	result := r.db.Model(&models.Location{}).
		Where("agent_id = ?", agentID).
		Count(&count)
	
	if result.Error != nil {
		return 0, result.Error
	}
	
	return count, nil
}