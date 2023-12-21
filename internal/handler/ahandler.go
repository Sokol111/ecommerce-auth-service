package handler

import (
	"errors"
	"github.com/Sokol111/ecommerce-auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) BindRoutes(engine *gin.Engine) {
	group := engine.Group("/auth")
	group.POST("/login", h.login)
	group.GET("/current", h.currentUser)
}

func (h *AuthHandler) login(c *gin.Context) {
	var request loginRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if token, err := h.service.Login(c, request.Login, request.Password); err == nil {
		c.JSON(200, loginResponse{token})
	} else {
		c.Error(err)
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
	}
}

func (h *AuthHandler) currentUser(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithError(http.StatusForbidden, errors.New("no Authorization header provided"))
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth || token == "" {
		c.AbortWithError(http.StatusForbidden, errors.New("couldn't find bearer token in Authorization header"))
		return
	}

	if u, err := h.service.GetUserByToken(c, token); err == nil {
		c.JSON(200, toJson(u))
	} else {
		c.Error(err)
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
	}
}
