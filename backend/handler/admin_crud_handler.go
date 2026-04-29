package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/reihanboo/kilas-admin/service"
)

type AdminCRUDHandler struct {
	service *service.AdminCRUDService
}

func NewAdminCRUDHandler(s *service.AdminCRUDService) *AdminCRUDHandler {
	return &AdminCRUDHandler{service: s}
}

func (h *AdminCRUDHandler) Summary(c *gin.Context) {
	data, err := h.service.Summary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminCRUDHandler) List(c *gin.Context) {
	entity := c.Param("entity")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	q := c.Query("q")

	data, err := h.service.List(entity, service.ListOptions{
		Q:      q,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminCRUDHandler) Get(c *gin.Context) {
	entity := c.Param("entity")
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.service.Get(entity, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminCRUDHandler) Create(c *gin.Context) {
	entity := c.Param("entity")
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.Create(entity, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, data)
}

func (h *AdminCRUDHandler) Update(c *gin.Context) {
	entity := c.Param("entity")
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.Update(entity, id, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AdminCRUDHandler) Delete(c *gin.Context) {
	entity := c.Param("entity")
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(entity, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func parseID(raw string) (uint, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}
