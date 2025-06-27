package controller

import (
	"net/http"
	"strconv"

	"my-gin-app/database/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOtherServices(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var services []model.OtherService
		if err := db.Find(&services).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)
	}
}

func GetOtherServiceByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var service model.OtherService
		if err := db.First(&service, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OtherService not found"})
			return
		}
		c.JSON(http.StatusOK, service)
	}
}

func CreateOtherService(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var service model.OtherService
		if err := c.ShouldBindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&service).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, service)
	}
}

func UpdateOtherService(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var service model.OtherService
		if err := db.First(&service, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OtherService not found"})
			return
		}
		if err := c.ShouldBindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&service).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, service)
	}
}

func DeleteOtherService(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		if err := db.Delete(&model.OtherService{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OtherService deleted"})
	}
}
