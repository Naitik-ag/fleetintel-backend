package repository

import (
	"errors"

	"gorm.io/gorm"
	"github.com/Naitik-ag/fleetintel-backend/internal/models"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

func (r *AgentRepository) Create(agent *models.DeliveryAgent) error {
	var existing models.DeliveryAgent
	result := r.db.Where("id = ?", agent.ID).First(&existing)
	
	if result.Error == nil {
		return errors.New("agent with this ID already exists")
	}
	
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	
	result = r.db.Where("phone = ?", agent.Phone).First(&existing)
	
	if result.Error == nil {
		return errors.New("agent with this phone number already exists")
	}
	
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}
	
	result = r.db.Create(agent)
	return result.Error
}

func (r *AgentRepository) FindByID(agentID string) (*models.DeliveryAgent, error) {
	var agent models.DeliveryAgent
	
	result := r.db.Where("id = ?", agentID).First(&agent)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &agent, nil
}

func (r *AgentRepository) FindAll(limit, offset int, status string) ([]models.DeliveryAgent, int64, error) {
	var agents []models.DeliveryAgent
	var totalCount int64
	
	query := r.db.Model(&models.DeliveryAgent{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&agents).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return agents, totalCount, nil
}

func (r *AgentRepository) Update(agentID string, updates map[string]interface{}) error {
	result := r.db.Model(&models.DeliveryAgent{}).
		Where("id = ?", agentID).
		Updates(updates)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("agent not found")
	}
	
	return nil
}

func (r *AgentRepository) UpdateStatus(agentID string, status string) error {
	validStatuses := map[string]bool{
		"available": true,
		"busy":      true,
		"offline":   true,
	}
	
	if !validStatuses[status] {
		return errors.New("invalid status. Must be: available, busy, or offline")
	}
	
	result := r.db.Model(&models.DeliveryAgent{}).
		Where("id = ?", agentID).
		Update("status", status)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("agent not found")
	}
	
	return nil
}

func (r *AgentRepository) Delete(agentID string) error {
	result := r.db.Where("id = ?", agentID).Delete(&models.DeliveryAgent{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("agent not found")
	}
	
	return nil
}

func (r *AgentRepository) SoftDelete(agentID string) error {
	return r.Update(agentID, map[string]interface{}{
		"is_active": false,
		"status":    "offline",
	})
}


func (r *AgentRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}
	
	err := r.db.Model(&models.DeliveryAgent{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}
	
	return counts, nil
}

func (r *AgentRepository) FindByPhone(phone string) (*models.DeliveryAgent, error) {
	var agent models.DeliveryAgent
	
	result := r.db.Where("phone = ?", phone).First(&agent)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &agent, nil
}