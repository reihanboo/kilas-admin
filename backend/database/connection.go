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

var AdminDB *gorm.DB
var KilasDB *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")

	adminDBName := "kilas_admin_db"
	if customName := os.Getenv("DATABASE_NAME"); customName != "" && customName != adminDBName {
		log.Fatal("DATABASE_NAME must be 'kilas_admin_db' for admin auth database")
	}

	kilasDBName := os.Getenv("KILAS_DATABASE_NAME")
	if kilasDBName == "" {
		kilasDBName = "kilas_db"
	}

	adminDB, err := openDB(dbUser, dbPass, dbHost, dbPort, adminDBName)
	if err != nil {
		log.Fatal("Error connecting to admin database: ", err)
	}
	kilasDB, err := openDB(dbUser, dbPass, dbHost, dbPort, kilasDBName)
	if err != nil {
		log.Fatal("Error connecting to kilas database: ", err)
	}

	AdminDB = adminDB
	KilasDB = kilasDB

	autoMigrateAdminDB()
	seedDefaultAdmin()
}

func openDB(username, password, host, port, dbName string) (*gorm.DB, error) {
	dsn := username + ":" + password + "@(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func autoMigrateAdminDB() {
	err := AdminDB.AutoMigrate(
		&model.User{},
		&model.Issue{},
	)
	if err != nil {
		log.Fatal("Error migrating admin database: ", err)
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
	err := AdminDB.Where("email = ?", email).First(&existing).Error
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

	if err := AdminDB.Create(&admin).Error; err != nil {
		log.Println("Error creating seeded admin user:", err)
	}
}
