package handlers

import (
	"insightly/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshTokensHandler interface {
	Refresh(c *gin.Context)
}
type RefreshTokensHandlerImpl struct {
	Rtr    services.RefreshTokensService
	Us     services.UserService
	Secret string
}

func NewRefreshTokensHandlerImpl(rtr services.RefreshTokensService, secret string, us services.UserService) *RefreshTokensHandlerImpl {
	return &RefreshTokensHandlerImpl{Rtr: rtr, Secret: secret, Us: us}
}

func (h *RefreshTokensHandlerImpl) Refresh(c *gin.Context) {
	var client struct {
		Token string `json:"refresh_token"`
	}

	//Читаем токен
	err := c.ShouldBindJSON(&client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Получаем токен
	t, err := h.Rtr.GetTokenByValue(client.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Достаем Id юзера
	token, err := h.Us.GenerateAccessToken(t.UserId, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token":token})

}
