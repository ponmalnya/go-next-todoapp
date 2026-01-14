package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     string
	DBPath   string
	FrontURL string
}

type Todo struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func SetupCORS(r *gin.Engine, frontURL string) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{frontURL},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))
}

func initDB(dbPath string) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("faild to connect database" + err.Error())
	}
	db.AutoMigrate(&Todo{})
	return db
}

func createTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&todo)
		c.JSON(http.StatusCreated, todo)
	}
}

func updateTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invaild ID"})
			return
		}
		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		var input Todo
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		todo.Title = input.Title
		todo.Completed = input.Completed
		db.Save(&todo)
		c.JSON(http.StatusOK, todo)

	}
}
func deleteTodoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		db.Delete(&todo)
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
	}
}

func getTodosHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		c.JSON(http.StatusOK, todos)
	}
}

func main() {

	godotenv.Load()

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	cfg := &Config{

		Port:     getEnv("PORT", "8080"),
		DBPath:   getEnv("DB_PATH", "todo.db"),
		FrontURL: getEnv("FRONT_URL", "http://localhost:3000"),
	}

	db := initDB(cfg.DBPath)
	r := gin.Default()
	SetupCORS(r, cfg.FrontURL)

	r.GET("/todos", getTodosHandler(db))
	r.POST("/todos", createTodoHandler(db))
	r.PUT("/todos/:id", updateTodoHandler(db))
	r.DELETE("/todos/:id", deleteTodoHandler(db))

	r.Run(":" + cfg.Port)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
