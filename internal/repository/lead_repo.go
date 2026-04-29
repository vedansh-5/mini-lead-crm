// This layer handles all direct interactions with PostgreSQL using GORM.
package repository

import (
	"mini-lead-crm/internal/models"

	"gorm.io/gorm"
)

type LeadRepository struct {
	db *gorm.DB
}

// create new instance of LeadRepository
func NewLeadRepository(db *gorm.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

// create & inserts a single lead in db
func (r *LeadRepository) Create(lead *models.Lead) error {
	return r.db.Create(lead).Error
}

func (r *LeadRepository) FindByID(id string) (*models.Lead, error) {
	var lead models.Lead
	err := r.db.First(&lead, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// retrives all leads, can filter by status
func (r *LeadRepository) FindAll(status string) ([]models.Lead, error) {
	var leads []models.Lead
	query := r.db

	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&leads).Error
	return leads, err
}

// update lead
func (r *LeadRepository) Update(lead *models.Lead) error {
	return r.db.Save(lead).Error
}

// delete lead
func (r *LeadRepository) Delete(id string) error {
	return r.db.Delete(&models.Lead{}, "id=?", id).Error
}
