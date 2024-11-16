package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// HealthController provides handlers for health checks
type HealthController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// NewHealthController initializes a new HealthController
func NewHealthController(db *gorm.DB, redisClient *redis.Client) *HealthController {
	return &HealthController{
		DB:          db,
		RedisClient: redisClient,
	}
}

// RegisterHealthRoutes sets up the health check routes
func (hc *HealthController) RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", hc.HealthCheck)
	router.GET("/health/db", hc.DatabaseHealthCheck)
	router.GET("/health/redis", hc.RedisHealthCheck)
}

// HealthCheck returns the overall service health
func (hc *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Booking Service is running",
	})
}

// DatabaseHealthCheck verifies database connectivity
func (hc *HealthController) DatabaseHealthCheck(c *gin.Context) {
	sqlDB, err := hc.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Database is down",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Database is healthy",
	})
}

// RedisHealthCheck verifies Redis connectivity
func (hc *HealthController) RedisHealthCheck(c *gin.Context) {
	if err := hc.RedisClient.Ping(c).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Redis is down",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Redis is healthy",
	})
}
