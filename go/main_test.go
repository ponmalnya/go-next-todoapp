package main
import (
	"bites"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
		
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
)

func setupTestDB() *gorm.DB{
	db, _ gorm.Open(sqlite.Open("file:: memory:?cache=shared"),gorm.Config{})
	db.AutoMigrate(&Todo{})
	return db
}