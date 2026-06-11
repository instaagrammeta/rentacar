package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/middleware"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// Login authenticates a user and returns access + refresh tokens.
func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	_ = c.ShouldBindJSON(&body)
	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Введите имя пользователя и пароль"})
		return
	}
	result, err := h.Svc.Authenticate(body.Username, body.Password, c.ClientIP())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Me returns the currently authenticated user.
func (h *Handler) Me(c *gin.Context) {
	user, err := h.Svc.GetUserByID(middleware.UserIDOf(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}
	c.JSON(http.StatusOK, user.ToMap())
}

// Logout revokes the current access token (Redis blacklist).
func (h *Handler) Logout(c *gin.Context) {
	if v, ok := c.Get(middleware.CtxJTI); ok {
		if jti, ok := v.(string); ok {
			h.Svc.Logout(jti, time.Now().Add(h.Cfg.JWTAccessExpiry))
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Вы вышли из системы"})
}

// Roles lists the available user roles.
func (h *Handler) Roles(c *gin.Context) {
	roles := make([]gin.H, 0)
	for _, r := range models.AllUserRoles() {
		roles = append(roles, gin.H{"value": string(r), "label": r.Label()})
	}
	c.JSON(http.StatusOK, roles)
}

// ListUsers returns all employee accounts (admin only).
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.Svc.ListUsers()
	if err != nil {
		respondError(c, err)
		return
	}
	out := make([]map[string]interface{}, 0, len(users))
	for i := range users {
		out = append(out, users[i].ToMap())
	}
	c.JSON(http.StatusOK, out)
}

// CreateUser creates an employee account (admin only).
func (h *Handler) CreateUser(c *gin.Context) {
	var in service.CreateUserInput
	_ = c.ShouldBindJSON(&in)
	user, err := h.Svc.CreateUser(in, actorFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user.ToMap())
}

// ChangePassword updates a user's password (admin only).
func (h *Handler) ChangePassword(c *gin.Context) {
	id, ok := paramID(c, "id")
	if !ok {
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.Svc.ChangePassword(id, body.Password); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пароль обновлён"})
}
