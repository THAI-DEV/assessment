package handler

import (
	"net/http"
	"strconv"

	"github.com/THAI-DEV/assessment/database"
	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

func Create(c *gin.Context) {
	var expense database.Expense
	err := c.ShouldBindJSON(&expense)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	id, err := database.CreateData(expense)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	if id > 0 {
		expense.Id = strconv.Itoa(id)
	}

	c.JSON(http.StatusCreated, expense)
}

func Read(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	result, err := database.ReadData(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
