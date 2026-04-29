// This file maps HTTP routes (like POST /leads) to Go functions using the Gin framework.
package handlers

import (
	"encoding/json"
	"mini-lead-crm/internal/models"
	"mini-lead-crm/internal/repository"
	"mini-lead-crm/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeadHandler struct {
	service *services.LeadService
	repo    *repository.LeadRepository
}

func NewLeadHandler(service *services.LeadService, repo *repository.LeadRepository) *LeadHandler {
	return &LeadHandler{service: service, repo: repo}
}

// POST /leads
func (h *LeadHandler) CreateLead(c *gin.Context) {
	var lead models.Lead

	// ShouldBindJSON validates againts tags (binding:required)
	if err := c.ShouldBindJSON(&lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateLead(&lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lead"})
		return
	}
	c.JSON(http.StatusCreated, lead)
}

// GET /leads/:id
func (h *LeadHandler) GetLeadByID(c *gin.Context) {
	id := c.Param("id")

	lead, err := h.service.GetLeadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}
	c.JSON(http.StatusOK, lead)
}

// GET /leads with optional ?status= filtering
func (h *LeadHandler) GetLeads(c *gin.Context) {
	status := c.Query("status") // Grab the query parameter

	leads, err := h.repo.FindAll(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}
	c.JSON(http.StatusOK, leads)
}

// PATCH /leads/:id/status
func (h *LeadHandler) UpdateLeadStatus(c *gin.Context) {
	id := c.Param("id")

	var payload struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	err := h.service.UpdateLeadStatus(id, models.LeadStatus(payload.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// PUT /leads/:id
func (h *LeadHandler) UpdateLead(c *gin.Context) {
	id := c.Param("id")
	var lead models.Lead

	if err := c.ShouldBindJSON(&lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateLead(id, &lead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lead"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead updated successfully"})
}

// DELETE /leads/:id
func (h *LeadHandler) DeleteLead(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteLead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lead"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead deleted successfully"})
}

// POST /leads/bulk
func (h *LeadHandler) BulkCreateLeads(c *gin.Context) {
	var leads []models.Lead

	if err := json.NewDecoder(c.Request.Body).Decode(&leads); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid array payload"})
		return
	}

	result := h.service.BulkCreate(leads)
	c.JSON(http.StatusMultiStatus, result)
}

// PUT /leads/bulk
func (h *LeadHandler) BulkUpdateLeads(c *gin.Context) {
	var leads []models.Lead

	if err := json.NewDecoder(c.Request.Body).Decode(&leads); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid array payload"})
		return
	}
	result := h.service.BulkUpdate(leads)
	c.JSON(http.StatusMultiStatus, result)
}
