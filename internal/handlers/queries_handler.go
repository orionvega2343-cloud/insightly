package handlers

import (
	"insightly/internal/middlewares"
	"insightly/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueriesHandler interface {
	CreateQueries(c *gin.Context)
	GetQueriesByUserId(c *gin.Context)
}

type QueriesHandlerImpl struct {
	Q services.QueriesService
}

func NewQueriesHandler(q services.QueriesService) *QueriesHandlerImpl {
	return &QueriesHandlerImpl{Q: q}
}

func (q *QueriesHandlerImpl) CreateQueries(c *gin.Context) {
	var queries struct {
		FileId   int    `json:"file_id"`
		Question string `json:"question"`
	}

	err := c.ShouldBindJSON(&queries)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, ok := middlewares.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	query, err := q.Q.CreateQueries(queries.FileId, id, queries.Question)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": query})
}

func (q *QueriesHandlerImpl) GetQueriesByUserId(c *gin.Context) {
	id, ok := middlewares.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	query, err := q.Q.GetQueriesByUserId(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": query})
}
