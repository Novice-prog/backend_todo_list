package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"time"
	"todo_list/internal/models"
)

var DB *gorm.DB

func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "./todos.db"
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get raw db: %w", err)
	}
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		return fmt.Errorf("auto migrate users: %v", err)
	}

	DB = db
	return nil
}

//------User----------

func CreateUser(input models.RegisterInput) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	if err := DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := DB.Where("email = ?", email).First(&u).Error
	return u, err
}

func GetUserByUsername(username string) (models.User, error) {
	var u models.User
	err := DB.Where("username = ?", username).First(&u).Error
	return u, err
}

func GetUserByID(id uint) (models.User, error) {
	var u models.User
	err := DB.First(&u, id).Error
	return u, err
}

// ----------todo-------------
func GetTodos(userID uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&todos).Error
	return todos, err
}

func CreateTodo(userID uint, input models.CreateTodoInput) (models.Todo, error) {
	todo := models.Todo{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
	}
	if err := DB.Create(&todo).Error; err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func UpdateTodo(userID uint, todoID uint, in models.UpdateTodoInput) error {
	var todo models.Todo
	if err := DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		return err
	}
	if in.Title != nil {
		todo.Title = *in.Title
	}
	if in.Description != nil {
		todo.Description = *in.Description
	}
	if in.Completed != nil {
		todo.Completed = *in.Completed
	}
	return DB.Save(&todo).Error
}

func DeleteTodo(userID, todoID uint) error {
	return DB.Where("id = ? AND user_id = ?", todoID, userID).
		Delete(&models.Todo{}).Error
}

func GetTodoByID(userID, todoID uint) (models.Todo, error) {
	var todo models.Todo
	err := DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	return todo, err
}
