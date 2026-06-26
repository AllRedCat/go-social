// Package posts handle posts from users
package posts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler sctructure
type Handler struct {
	service Service
}

// NewHandler -> new posts handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	postsGroup := router.Group("/api/posts")
	{
		postsGroup.POST("/post/:id", h.Post)
	}
}

func (h *Handler) Post(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req PostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	postResp, err := h.service.Post(c.Request.Context(), req, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "post": postResp})
}
