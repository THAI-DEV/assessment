package handler

import (
	"net/http"
	"strconv"

	"github.com/THAI-DEV/assessment/database"
	"github.com/gin-gonic/gin"
)

type ExpenseBody struct {
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

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

func ReadOne(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	result, err := database.ReadOneData(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func Update(c *gin.Context) {
	idParam := c.Param("id")
	// iid, _ := strconv.Atoi(idParam)

	var expenseBody ExpenseBody

	err := c.ShouldBindJSON(&expenseBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	expense := database.Expense{
		Id:     idParam,
		Title:  expenseBody.Title,
		Amount: expenseBody.Amount,
		Note:   expenseBody.Note,
		Tags:   expenseBody.Tags,
	}

	id, err := database.UpdateData(expense)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	if id > 0 {
		expense.Id = strconv.Itoa(id)
	}

	c.JSON(http.StatusOK, expense)
}

func ReadAll(c *gin.Context) {
	result, err := database.ReadAllData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Request header authorization is empty",
			})
			c.Abort()
			return
		}

		if authHeader != "November 10, 2009" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorization",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
