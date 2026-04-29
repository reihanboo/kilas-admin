package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/reihanboo/kilas-admin/dto"
	"github.com/reihanboo/kilas-admin/middleware"
	"github.com/reihanboo/kilas-admin/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: token,
		User: dto.CurrentUser{
			ID:    user.ID,
			Name:  user.Username,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.Me(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, dto.CurrentUser{
		ID:    user.ID,
		Name:  user.Username,
		Email: user.Email,
		Role:  user.Role,
	})
}

type IssueHandler struct {
	issueService *service.IssueService
}

func NewIssueHandler(issueService *service.IssueService) *IssueHandler {
	return &IssueHandler{issueService: issueService}
}

func (h *IssueHandler) CreateIssue(c *gin.Context) {
	var req dto.CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue, err := h.issueService.CreateIssue(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create issue"})
		return
	}

	c.JSON(http.StatusCreated, issue)
}

func (h *IssueHandler) ListIssues(c *gin.Context) {
	status := c.Query("status")
	issues, err := h.issueService.ListIssues(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch issues"})
		return
	}
	c.JSON(http.StatusOK, issues)
}

func (h *IssueHandler) GetIssue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	issue, err := h.issueService.GetIssueByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
		return
	}
	c.JSON(http.StatusOK, issue)
}

func (h *IssueHandler) UpdateIssue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue, err := h.issueService.UpdateIssue(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update issue"})
		return
	}

	c.JSON(http.StatusOK, issue)
}

func (h *IssueHandler) DashboardSummary(c *gin.Context) {
	summary, err := h.issueService.DashboardSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
		return
	}
	c.JSON(http.StatusOK, summary)
}
