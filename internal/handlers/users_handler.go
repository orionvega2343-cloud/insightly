package handlers

import (
	"insightly/internal/models"
	"insightly/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type UserHandlerImpl struct {
	Us     services.UserService
	Rts    services.RefreshTokensService
	Secret string
}

func NewUserHandler(us services.UserService, rts services.RefreshTokensService, secret string) *UserHandlerImpl {
	return &UserHandlerImpl{Us: us, Rts: rts, Secret: secret}
}

func (h *UserHandlerImpl) Register(c *gin.Context) {
	var u models.User

	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Us.Register(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
	return
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	var u models.User

	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, token, err := h.Us.Login(u.Email, u.Password, h.Secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": token})
}
