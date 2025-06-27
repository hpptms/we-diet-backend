package controller

import (
	"net/http"
	"strconv"

	"my-gin-app/database/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPermissions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var permissions []model.Permission
		if err := db.Find(&permissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}

func GetPermissionByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var permission model.Permission
		if err := db.First(&permission, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusOK, permission)
	}
}

func CreatePermission(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var permission model.Permission
		if err := c.ShouldBindJSON(&permission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&permission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, permission)
	}
}

func UpdatePermission(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var permission model.Permission
		if err := db.First(&permission, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		if err := c.ShouldBindJSON(&permission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&permission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permission)
	}
}

func DeletePermission(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		if err := db.Delete(&model.Permission{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Permission deleted"})
	}
}
