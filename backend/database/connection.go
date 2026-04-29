package database

import (
	"errors"
	"log"
	"os"

	"github.com/reihanboo/kilas-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "kilas_admin_db"
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")

	dsn := dbUser + ":" + dbPass + "@(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	DB = db
	autoMigrate()
	seedDefaultAdmin()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Transaction{},
		&model.Deck{},
		&model.Card{},
		&model.AIGenerationHistory{},
		&model.Issue{},
	)
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}
}

func seedDefaultAdmin() {
	email := os.Getenv("ADMIN_SEED_EMAIL")
	password := os.Getenv("ADMIN_SEED_PASSWORD")
	username := os.Getenv("ADMIN_SEED_NAME")

	if email == "" || password == "" {
		return
	}
	if username == "" {
		username = "Kilas Admin"
	}

	var existing model.User
	err := DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error checking seeded admin user:", err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing seed admin password:", err)
		return
	}

	admin := model.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
		Role:     "admin",
		Provider: "local",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Println("Error creating seeded admin user:", err)
	}
}
