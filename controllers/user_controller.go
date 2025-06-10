package controllers

import (
	"net/http"

	"github.com/Run-Tu/go-scaffold/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

type UserInput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (ctrl *UserController) GetALLUsers(c *gin.Context) {
	var users []models.User
	ctrl.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
	}

	var user models.User
	if err := ctrl.DB.First(&user, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 客户端未提供ID,则自动生成
	if input.ID == "" {
		input.ID = uuid.NewString()
	}

	user := models.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}

	ctrl.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{"data": user})

}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User
	if err := ctrl.DB.First(&user, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	
	ctrl.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
