package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/reihanboo/kilas-admin/database"
	"github.com/reihanboo/kilas-admin/handler"
	"github.com/reihanboo/kilas-admin/router"
	"github.com/reihanboo/kilas-admin/service"
)

func main() {
	_ = godotenv.Load()

	appAddress := os.Getenv("APP_ADDRESS")
	if appAddress == "" {
		appAddress = "0.0.0.0"
	}
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8081"
	}

	database.InitDB()

	authService := service.NewAuthService(database.AdminDB)
	issueService := service.NewIssueService(database.AdminDB)
	adminService := service.NewAdminCRUDService(database.KilasDB, database.AdminDB)
	authHandler := handler.NewAuthHandler(authService)
	issueHandler := handler.NewIssueHandler(issueService)
	adminHandler := handler.NewAdminCRUDHandler(adminService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.RegisterRoutes(r, authHandler, issueHandler, adminHandler)

	log.Printf("kilas-admin backend running on %s:%s", appAddress, appPort)
	if err := r.Run(appAddress + ":" + appPort); err != nil {
		log.Fatal(err)
	}
}
