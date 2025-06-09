package controllers

import (
	"net/http"

	"github.com/Run-Tu/go-scaffold/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (ctrl *UserController) GetALLUsers(c *gin.Context) {
	var users []models.User
	ctrl.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}
