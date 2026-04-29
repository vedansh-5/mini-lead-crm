package services

import (
	"fmt"
	"mini-lead-crm/internal/cache"
	"mini-lead-crm/internal/models"
	"mini-lead-crm/internal/repository"

	"github.com/gin-gonic/gin/binding"
)

type LeadService struct {
	repo  *repository.LeadRepository
	cache *cache.RedisClient
}

func NewLeadService(repo *repository.LeadRepository, cache *cache.RedisClient) *LeadService {
	return &LeadService{repo: repo, cache: cache}
}

func (s *LeadService) CreateLead(lead *models.Lead) error {
	lead.Status = models.StatusNew
	return s.repo.Create(lead)
}

func (s *LeadService) GetLeadByID(id string) (*models.Lead, error) {
	// try cache for lead
	cachedLead, err := s.cache.Get(id)
	if err == nil && cachedLead != nil {
		return cachedLead, nil
	}

	lead, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	s.cache.Set(lead)
	return lead, nil
}

func (s *LeadService) UpdateLeadStatus(id string, newStatus models.LeadStatus) error {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	//check if transition allowed
	if !models.IsValidTransition(lead.Status, newStatus) {
		return fmt.Errorf("invalud status transition from %s to %s", lead.Status, newStatus)
	}

	lead.Status = newStatus
	err = s.repo.Update(lead)
	if err == nil {
		s.cache.Invalidate(id)
	}
	return err
}

// delete lead by ID
func (s *LeadService) DeleteLead(id string) error {
	err := s.repo.Delete(id)
	if err == nil {
		s.cache.Invalidate(id)
	}
	return err
}

// updates a lead's fields
func (s *LeadService) UpdateLead(id string, updatedLead *models.Lead) error {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	lead.Name = updatedLead.Name
	lead.Email = updatedLead.Email
	lead.Phone = updatedLead.Phone
	lead.Source = updatedLead.Source

	err = s.repo.Update(lead)
	if err == nil {
		s.cache.Invalidate(id) // remove cache since data changed
	}
	return err
}

// create multiple leads, report success and failures
func (s *LeadService) BulkCreate(leads []models.Lead) map[string]interface{} {
	successful := 0
	failed := 0
	var results []map[string]interface{}

	for i, lead := range leads {
		// validate the individual lead
		if err := binding.Validator.ValidateStruct(&lead); err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"index":   i,
				"success": false,
				"error":   err.Error(),
			})
			continue // Skip creating invalid lead
		}
		err := s.CreateLead(&lead)

		if err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"index":   i,
				"success": false,
				"error":   err.Error(),
			})
		} else {
			successful++
			results = append(results, map[string]interface{}{
				"index":   i,
				"success": true,
				"lead":    lead,
			})
		}
	}

	return map[string]interface{}{
		"total":      len(leads),
		"successful": successful,
		"failed":     failed,
		"results":    results,
	}
}

// BulkUpdate updates multiple leads
func (s *LeadService) BulkUpdate(leads []models.Lead) map[string]interface{} {
	successful := 0
	failed := 0
	var results []map[string]interface{}
	for i, lead := range leads {
		// validate the individual lead
		if err := binding.Validator.ValidateStruct(&lead); err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"index":   i,
				"success": false,
				"error":   err.Error(),
			})
			continue // Skip creating invalid lead
		}
		err := s.UpdateLead(lead.ID.String(), &lead)
		if err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"id":      lead.ID,
				"success": false,
				"error":   err.Error(),
			})
		} else {
			successful++
			results = append(results, map[string]interface{}{
				"id":      lead.ID,
				"success": true,
			})
		}
	}
	return map[string]interface{}{
		"total":      len(leads),
		"successful": successful,
		"failed":     failed,
		"results":    results,
	}
}
