// Package handler implements the Gin HTTP handlers, mapping the REST API
// contract of the original Flask backend onto the Go service layer.
package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/config"
	"github.com/instaagrammeta/rentacar/backend-go/internal/middleware"
	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

// Handler holds the dependencies shared by every HTTP handler.
type Handler struct {
	Svc  *service.Service
	Auth *middleware.AuthService
	Cfg  *config.Config
}

// New builds a Handler.
func New(svc *service.Service, auth *middleware.AuthService, cfg *config.Config) *Handler {
	return &Handler{Svc: svc, Auth: auth, Cfg: cfg}
}

// actorFrom extracts the authenticated actor from the request context.
func actorFrom(c *gin.Context) *service.Actor {
	id := middleware.UserIDOf(c)
	actor := &service.Actor{Username: middleware.UsernameOf(c), IP: c.ClientIP()}
	if id > 0 {
		actor.ID = &id
	}
	return actor
}

// respondError translates an error into a consistent JSON response.
func respondError(c *gin.Context, err error) {
	var appErr *apperr.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}
	log.Printf("internal error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
}

// pagination parses page / per_page query params (page>=1, 1<=per_page<=200).
func pagination(c *gin.Context) (page, perPage, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ = strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 200 {
		perPage = 200
	}
	return page, perPage, (page - 1) * perPage
}

// paramID parses an unsigned integer path parameter.
func paramID(c *gin.Context, name string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ресурс не найден"})
		return 0, false
	}
	return uint(id), true
}

// queryUint parses an optional unsigned integer query parameter (0 if absent).
func queryUint(c *gin.Context, name string) uint {
	id, _ := strconv.ParseUint(c.Query(name), 10, 64)
	return uint(id)
}
