package handlers

import (
	"insightly/internal/middlewares"
	"insightly/internal/services"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FilesHandler interface {
	CreateFiles(c *gin.Context)
	GetFilesByUserId(c *gin.Context)
	DeleteFile(c *gin.Context)
}

type FilesHandlerImpl struct {
	F services.FilesService
}

func NewFilesHandler(f services.FilesService) *FilesHandlerImpl {
	return &FilesHandlerImpl{F: f}
}

// CreateFiles godoc
// @Summary Загрузка CSV файла
// @Tags files
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV файл"
// @Success 200 {object} map[string]interface{} "data: models.Files"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /files/upload [post]
func (h *FilesHandlerImpl) CreateFiles(c *gin.Context) {
	//Получаем файл
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, ok := middlewares.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//Берем байты для передачи в сервис
	open, err := f.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer open.Close()

	read, err := io.ReadAll(open)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := h.F.CreateFiles(id, f.Filename, read)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": file})
}

// GetFilesByUserId godoc
// @Summary Получение списка файлов текущего пользователя
// @Tags files
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "data: []models.Files"
// @Failure 401 {object} map[string]string
// @Router /files [get]
func (h *FilesHandlerImpl) GetFilesByUserId(c *gin.Context) {
	id, ok := middlewares.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	files, err := h.F.GetFilesByUserId(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

// DeleteFile godoc
// @Summary Удаление файла по ID
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID файла"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /files/{id} [delete]
func (h *FilesHandlerImpl) DeleteFile(c *gin.Context) {
	fileIdParam := c.Param("id")

	fileId, err := strconv.Atoi(fileIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, ok := middlewares.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := h.F.GetFileById(fileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file.UserId != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = h.F.DeleteFile(fileId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
