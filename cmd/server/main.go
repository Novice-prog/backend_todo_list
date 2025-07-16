package main

import (
	"log"
	"os"
	"time"

	"todo_list/internal/database"
	"todo_list/internal/handlers"
	"todo_list/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("db init: %v", err)
	}

	/* ── 2. Создаём Gin-роутер ── */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery()) // базовые middleware

	/* ── 3. Лояльный CORS (чтобы Postman/фронт могли дергать API) ── */
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // разрешаем всё; сузьте при необходимости
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") }) // health-check
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	api := r.Group("/api", middleware.AuthMiddleware())
	{
		api.GET("/profile", handlers.GetProfile)

		api.GET("/todos", handlers.GetTodos)
		api.POST("/todos", handlers.CreateTodo)
		api.PATCH("/todos/:id", handlers.UpdateTodo)        // частичное обновление
		api.PATCH("/todos/:id/toggle", handlers.ToggleTodo) // переключить completed
		api.DELETE("/todos/:id", handlers.DeleteTodo)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server start: %v", err)
	}
}
