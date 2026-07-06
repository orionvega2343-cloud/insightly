package handlers

import (
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

func (h *FilesHandlerImpl) CreateFiles(c *gin.Context) {
	//Получаем файл
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Приводим к assert type 2 проверками
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := userId.(int)
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

func (h *FilesHandlerImpl) GetFilesByUserId(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := userId.(int)
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

func (h *FilesHandlerImpl) DeleteFile(c *gin.Context) {
	fileIdParam := c.Param("id")

	fileId, err := strconv.Atoi(fileIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := userId.(int)
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
