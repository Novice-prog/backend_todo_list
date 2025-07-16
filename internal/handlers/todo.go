package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo_list/internal/database"
	"todo_list/internal/models"
)

func currentUserID(c *gin.Context) (uint, bool) {
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GET /api/todos
func GetTodos(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todos, err := database.GetTodos(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, todos)
}

// POST /api/todos
func CreateTodo(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	var in models.CreateTodoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := database.CreateTodo(uid, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, todo)
}

// PATCH /api/todos/:id
func UpdateTodo(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var in models.UpdateTodoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.UpdateTodo(uid, uint(todoID), in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "todo update"})
}

func ToggleTodo(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := database.GetTodoByID(uid, uint(todoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newState := !todo.Completed
	in := models.UpdateTodoInput{Completed: &newState}

	if err := database.UpdateTodo(uid, uint(todoID), in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	update, _ := database.GetTodoByID(uid, uint(todoID))
	c.JSON(http.StatusOK, update)
}

func DeleteTodo(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := database.DeleteTodo(uid, uint(todoID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}
