package router

import (
	"github.com/gin-gonic/gin"
	"github.com/reihanboo/kilas-admin/handler"
	"github.com/reihanboo/kilas-admin/middleware"
)

func RegisterRoutes(r *gin.Engine, authHandler *handler.AuthHandler, issueHandler *handler.IssueHandler, adminHandler *handler.AdminCRUDHandler) {
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)

	// Public issue report endpoint used by Kilas users
	api.POST("/issues/report", issueHandler.CreateIssue)

	// Generic CRUD for admin CMS
	admin := api.Group("/admin", middleware.AuthMiddleware())
	admin.GET("/summary", adminHandler.Summary)
	admin.GET("/:entity", adminHandler.List)
	admin.GET("/:entity/:id", adminHandler.Get)
	admin.POST("/:entity", adminHandler.Create)
	admin.PUT("/:entity/:id", adminHandler.Update)
	admin.DELETE("/:entity/:id", adminHandler.Delete)
}
