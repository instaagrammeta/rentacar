// Package middleware contains Gin middleware for authentication, authorisation
// and consistent error responses.
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/instaagrammeta/rentacar/backend-go/internal/auth"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// Context keys used to share the authenticated identity with handlers.
const (
	CtxUserID   = "userID"
	CtxRole     = "role"
	CtxUsername = "username"
	CtxJTI      = "jti"
)

// RolePermissions maps a logical module permission to the roles allowed to use
// it, mirroring ROLE_PERMISSIONS from the Flask backend.
var RolePermissions = map[string][]models.UserRole{
	"dashboard":    {models.RoleAdministrator, models.RoleRentalManager, models.RoleCashier, models.RoleOperator},
	"clients":      {models.RoleAdministrator, models.RoleRentalManager},
	"cars":         {models.RoleAdministrator, models.RoleRentalManager},
	"reservations": {models.RoleAdministrator, models.RoleRentalManager, models.RoleOperator},
	"rentals":      {models.RoleAdministrator, models.RoleRentalManager},
	"returns":      {models.RoleAdministrator, models.RoleRentalManager},
	"payments":     {models.RoleAdministrator, models.RoleCashier},
	"blacklist":    {models.RoleAdministrator, models.RoleRentalManager},
	"accidents":    {models.RoleAdministrator, models.RoleRentalManager},
	"reports":      {models.RoleAdministrator},
	"settings":     {models.RoleAdministrator},
	"users":        {models.RoleAdministrator},
	"backups":      {models.RoleAdministrator},
}

// AuthService bundles the dependencies the auth middleware needs.
type AuthService struct {
	Tokens *auth.Manager
	Redis  *redis.Client
}

// blacklistKey is the Redis key used to revoke a token by its JWT id.
func blacklistKey(jti string) string { return "jwt:blacklist:" + jti }

// RequireAuth validates the access token and stores the identity in context.
func (a *AuthService) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			return
		}

		claims, err := a.Tokens.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Срок действия токена истёк"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}
		if claims.Type != auth.AccessToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
			return
		}

		// Reject revoked tokens (logout) when Redis is available.
		if a.Redis != nil && claims.ID != "" {
			if n, _ := a.Redis.Exists(context.Background(), blacklistKey(claims.ID)).Result(); n > 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
				return
			}
		}

		userID, _ := strconv.ParseUint(claims.Subject, 10, 64)
		c.Set(CtxUserID, uint(userID))
		c.Set(CtxRole, claims.Role)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxJTI, claims.ID)
		c.Next()
	}
}

// RequirePermission enforces access based on the RolePermissions matrix.
func (a *AuthService) RequirePermission(permission string) gin.HandlerFunc {
	allowed := RolePermissions[permission]
	return func(c *gin.Context) {
		role := RoleOf(c)
		if !roleIn(role, allowed) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав для выполнения операции"})
			return
		}
		c.Next()
	}
}

// RequireRole enforces that the caller holds one of the given roles.
func (a *AuthService) RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !roleIn(RoleOf(c), roles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав для выполнения операции"})
			return
		}
		c.Next()
	}
}

// RoleOf returns the authenticated role from the context.
func RoleOf(c *gin.Context) models.UserRole {
	if v, ok := c.Get(CtxRole); ok {
		if role, ok := v.(models.UserRole); ok {
			return role
		}
	}
	return ""
}

// UserIDOf returns the authenticated user id from the context (0 if absent).
func UserIDOf(c *gin.Context) uint {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// UsernameOf returns the authenticated username from the context.
func UsernameOf(c *gin.Context) string {
	if v, ok := c.Get(CtxUsername); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func roleIn(role models.UserRole, set []models.UserRole) bool {
	for _, r := range set {
		if r == role {
			return true
		}
	}
	return false
}
