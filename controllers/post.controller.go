package controllers

import (
	"BlogApp/config"
	"BlogApp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// private routes

func CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	uidFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Invalid user ID type"})
		return
	}
	type BlogInput struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var input BlogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	blog := models.Blog{
		Title:   input.Title,
		Content: input.Content,
		UserID:  uint(uidFloat),
	}

	if err := config.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      blog.ID,
		"title":   blog.Title,
		"content": blog.Content,
	})
}

func UpdateById(c *gin.Context) {
	// Parse blog ID
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid blog ID"})
		return
	}

	// Get user ID from context
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	floatID, ok := rawID.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user ID type"})
		return
	}
	userID := uint(floatID)

	// Find the blog post
	var blog models.Blog
	if err := config.DB.Where("id = ? AND user_id = ?", blogID, userID).First(&blog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Post not found or unauthorized"})
		return
	}

	// Bind input
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Update fields
	blog.Title = input.Title
	blog.Content = input.Content

	if err := config.DB.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Post updated", "blog": blog})
}

func DeleteById(c *gin.Context) {
	// Parse blog ID from URL param
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid blog ID"})
		return
	}

	// Get user ID from context (JWT)
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	floatID, ok := rawID.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid user ID type"})
		return
	}
	userID := uint(floatID)

	// Find the blog and ensure ownership
	var blog models.Blog
	if err := config.DB.Where("id = ? AND user_id = ?", blogID, userID).First(&blog).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Post not found or unauthorized"})
		return
	}

	// Delete the blog post
	if err := config.DB.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Post deleted successfully"})
}

// public routes
func GetAllPosts(c *gin.Context) {
	// Default values
	page := 1
	limit := 10
	sort := c.DefaultQuery("sort", "desc")
	search := c.Query("search")
	userIDStr := c.Query("user_id")

	// Parse page and limit
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var blogs []models.Blog
	var total int64

	query := config.DB.Model(&models.Blog{})

	// Filter by search keyword
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Filter by user ID (converted to int)
	if userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			query = query.Where("user_id = ?", userID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user_id"})
			return
		}
	}

	// Count total posts after filters
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to count posts"})
		return
	}

	offset := (page - 1) * limit
	order := "created_at desc"
	if sort == "asc" {
		order = "created_at asc"
	}

	// Retrieve posts with pagination and sorting
	if err := query.Order(order).Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"posts": blogs,
	})
}

func GetPostById(c *gin.Context) {
	blogID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid post ID"})
		return
	}

	var blog models.Blog
	if err := config.DB.First(&blog, uint(blogID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, blog)

}
