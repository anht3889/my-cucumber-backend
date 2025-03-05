package api

import (
	"my-cucumber-backend/models"
	"my-cucumber-backend/services"

	"github.com/gin-gonic/gin"
)

func CreateChartHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var chart models.Chart
	if err := c.ShouldBindJSON(&chart); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	typedUser := user.(*models.User)
	chart.UserID = typedUser.ID
	if err := services.CreateChart(&chart); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, chart)
}

func GetChartsHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	typedUser := user.(*models.User)
	charts, err := services.GetChartsByUser(typedUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, charts)
}

func CreateDataTableHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var table models.DataTable
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	typedUser := user.(*models.User)
	table.UserID = typedUser.ID
	if err := services.CreateDataTable(&table); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, table)
}

func GetDataTablesHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	typedUser := user.(*models.User)
	tables, err := services.GetDataTablesByUser(typedUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tables)
}

// Similar handlers for DataTables...
