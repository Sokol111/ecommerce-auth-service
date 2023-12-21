package handler

import (
	"errors"
	"github.com/Sokol111/ecommerce-auth-service/internal/model"
	"github.com/Sokol111/ecommerce-auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Enabled  bool   `json:"enabled"`
}

type updateUserRequest struct {
	Id      string `json:"id" binding:"required"`
	Version int32  `json:"version" binding:"required"`
	Login   string `json:"login" binding:"required"`
	Enabled bool   `json:"enabled"`
}

type userResponse struct {
	Id      string `json:"id"`
	Version int32  `json:"version"`
	Login   string `json:"login"`
	Enabled bool   `json:"enabled"`
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) BindRoutes(engine *gin.Engine) {
	group := engine.Group("/user")
	group.POST("/create", h.createUser)
	group.PUT("/update", h.updateUser)
	group.GET("/list", h.getAllUsers)
	group.GET("/:id", h.getUserById)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var request createUserRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if u, err := h.service.Create(c, model.User{Login: request.Login, Enabled: request.Enabled}, request.Password); err == nil {
		c.JSON(200, toJson(u))
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h *UserHandler) getUserById(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("id parameter is empty"))
		return
	}

	if u, err := h.service.GetById(c, id); err == nil {
		c.JSON(200, toJson(u))
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h *UserHandler) getAllUsers(c *gin.Context) {
	if users, err := h.service.GetUsers(c); err == nil {
		response := make([]userResponse, 0, len(users))
		for _, u := range users {
			response = append(response, toJson(u))
		}
		c.JSON(200, response)
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h *UserHandler) updateUser(c *gin.Context) {
	var request updateUserRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if u, err := h.service.Update(c, model.User{ID: request.Id, Version: request.Version, Login: request.Login, Enabled: request.Enabled}); err == nil {
		c.JSON(200, toJson(u))
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func toJson(u model.User) userResponse {
	return userResponse{Id: u.ID, Version: u.Version, Login: u.Login, Enabled: u.Enabled}
}
