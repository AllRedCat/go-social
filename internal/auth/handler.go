package auth

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler structure
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes groups and registers all auth routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		
		// In a real application, these routes should be protected by a JWT middleware.
		// For now, we are passing the ID in the URL to identify the user.
		authGroup.POST("/avatar/:id", h.UploadAvatar)
		authGroup.PUT("/update/:id", h.UpdateUser)
		authGroup.DELETE("/delete/:id", h.SoftDelete)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Validate JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	// 2. Call service to hash password and save to database
	userResp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return success
	c.JSON(http.StatusCreated, gin.H{"message": "User successfully created", "user": userResp})
}

// Login handles user authentication and cookie generation
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Get JWT token from service
	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set JWT in HttpOnly cookie
	// Name, Value, MaxAge(Seconds), Path, Domain, Secure(true = HTTPS), HttpOnly(Prevents XSS)
	c.SetCookie("jwt_token", token, 3600*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// UploadAvatar handles the avatar image upload
func (h *Handler) UploadAvatar(c *gin.Context) {
	// Parse user ID from URL
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Extract the file from the form-data
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Generate a unique filename using timestamp and user ID
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(file.Filename))
	savePath := fmt.Sprintf("uploads/avatars/%s", filename)

	// Save the file physically
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Update the user's avatar URL in the database
	avatarURL := "/" + savePath
	if err := h.service.UpdateAvatar(c.Request.Context(), uint(userID), avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar successfully updated",
		"avatar_url": avatarURL,
	})
}

// UpdateUser handles the update of user information (name, email)
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	// Build the entity to pass to the service
	userToUpdate := &User{
		ID:    uint(userID),
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.service.UpdateUser(c.Request.Context(), userToUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully updated"})
}

// SoftDelete handles the logical deletion of a user
func (h *Handler) SoftDelete(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}
