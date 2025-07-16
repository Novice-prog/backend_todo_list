package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"todo_list/internal/database"
	"todo_list/internal/models"
)

const (
	cookieName = "token"
	cookieAge  = 24 * 60 * 60
)

var jwtSecret = func() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("change-me")
}()

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func setAuthCookie(c *gin.Context, token string) {
	secure := c.Request.TLS != nil // true, если соединение https
	c.SetCookie(
		cookieName,
		token,
		cookieAge,
		"/",
		"",     // domain — по умолчанию текущий
		secure, // Secure
		true,   // HttpOnly
	)
}

func Register(c *gin.Context) {
	var in models.RegisterInput
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.CreateUser(in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := generateToken(uint(user.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func Login(c *gin.Context) {
	var in models.LoginInput
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByUsername(in.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func GetProfile(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := database.GetUserByID(uid.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}
