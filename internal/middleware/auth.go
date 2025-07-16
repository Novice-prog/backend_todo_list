package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"todo_list/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = func() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	// dev-fallback – обязательно переопределите в проде!
	return []byte("change-me")
}()

func tokenFromRequest(c *gin.Context) string {
	if t, err := c.Cookie("token"); err == nil && t != "" {
		return t
	}
	if h := c.GetHeader("Authorization"); strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	return ""
}

func unauth(c *gin.Context, msg string) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
	} else {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}

/* ---------- основное middleware ---------- */

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) достаём токен
		tokenStr := tokenFromRequest(c)
		if tokenStr == "" {
			unauth(c, "token required")
			return
		}

		// 2) парсим и проверяем подпись
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			unauth(c, "invalid token")
			return
		}

		// 3) разбираем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauth(c, "invalid token claims")
			return
		}

		uidFloat, ok := claims["user_id"].(float64)
		if !ok {
			unauth(c, "invalid user_id claim")
			return
		}
		uid := uint(uidFloat)

		// 4) проверяем, что пользователь ещё существует (токен может быть «осиротевшим»)
		if _, err := database.GetUserByID(uid); err != nil {
			unauth(c, "user not found")
			return
		}

		// 5) всё ок – пишем user_id в контекст и пропускаем дальше
		c.Set("user_id", uid)
		c.Next()
	}
}
