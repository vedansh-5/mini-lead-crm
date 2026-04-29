package main

import (
	"log"
	"mini-lead-crm/internal/cache"
	"mini-lead-crm/internal/handlers"
	"mini-lead-crm/internal/models"
	"mini-lead-crm/internal/repository"
	"mini-lead-crm/internal/services"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=crm port=5432 sslmode=disable"
	}

	redisAddr :=

		os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "loaclhost:6379"
	}

	// connect to postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// run migrations to create tables
	err = db.AutoMigrate(&models.Lead{})
	if err != nil {

		log.Fatal("Failed to run migrations:", err)
	}

	// init components
	repo := repository.NewLeadRepository(db)
	redisCache := cache.NewRedisClient(redisAddr)
	service := services.NewLeadService(repo, redisCache)
	handler := handlers.NewLeadHandler(service, repo)

	router := gin.Default()

	router.POST("/leads", handler.CreateLead)
	router.GET("/leads", handler.GetLeads)
	router.GET("/leads/:id", handler.GetLeadByID)
	router.PATCH("/leads/:id/status", handler.UpdateLeadStatus)
	router.POST("/leads/bulk", handler.BulkCreateLeads)
	router.PUT("/leads/:id", handler.UpdateLead)
	router.DELETE("/leads/:id", handler.DeleteLead)
	router.PUT("/leads/bulk", handler.BulkUpdateLeads)

	router.Run(":8080")
}
