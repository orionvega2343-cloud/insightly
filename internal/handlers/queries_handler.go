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

// CreateQueries godoc
// @Summary Отправка вопроса по загруженному файлу для AI-анализа
// @Tags queries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "file_id, question"
// @Success 200 {object} map[string]interface{} "data: models.Queries"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /analyze [post]
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

// GetQueriesByUserId godoc
// @Summary Получение истории запросов текущего пользователя
// @Tags queries
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "data: []models.Queries"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /analyze/history [get]
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
