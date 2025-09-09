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
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Published bool   `json:"published"`
		Draft     bool   `json:"draft"`
	}
	var input BlogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// Enforce mutual exclusivity: if published is true, draft must be false, and vice versa
	published := input.Published
	draft := input.Draft
	if published {
		draft = false
	} else if draft {
		published = false
	}
	blog := models.Blog{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    uint(uidFloat),
		Published: published,
		Draft:     draft,
	}

	if err := config.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        blog.ID,
		"title":     blog.Title,
		"content":   blog.Content,
		"published": blog.Published,
		"draft":     blog.Draft,
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
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Published bool   `json:"published"`
		Draft     bool   `json:"draft"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Enforce mutual exclusivity: if published is true, draft must be false, and vice versa
	published := input.Published
	draft := input.Draft
	if published {
		draft = false
	} else if draft {
		published = false
	}
	blog.Title = input.Title
	blog.Content = input.Content
	blog.Published = published
	blog.Draft = draft

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

// func GetAllPosts(c *gin.Context) {
// 	// Default pagination and sorting values
// 	page := 1
// 	limit := 10
// 	sort := c.DefaultQuery("sort", "desc")
// 	search := c.Query("search")
// 	userIDStr := c.Query("user_id") // optional filter

// 	// Parse page and limit
// 	if p := c.Query("page"); p != "" {
// 		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
// 			page = parsedPage
// 		}
// 	}
// 	if l := c.Query("limit"); l != "" {
// 		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
// 			limit = parsedLimit
// 		}
// 	}

// 	var blogs []models.Blog
// 	var total int64

// 	query := config.DB.Model(&models.Blog{}).Where("published = ?", true)

// 	// Search filter
// 	if search != "" {
// 		query = query.Where("title LIKE ?", "%"+search+"%")
// 	}

// 	// User filter (My Posts)
// 	if userIDStr != "" {
// 		userID, err := strconv.Atoi(userIDStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user_id"})
// 			return
// 		}
// 		query = query.Where("user_id = ?", uint(userID))
// 	}

// 	// Count total posts after filters
// 	if err := query.Count(&total).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to count posts"})
// 		return
// 	}

// 	// Sorting and pagination
// 	offset := (page - 1) * limit
// 	order := "created_at desc"
// 	if sort == "asc" {
// 		order = "created_at asc"
// 	}

// 	if err := query.Order(order).Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve posts"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"page":  page,
//			"limit": limit,
//			"total": total,
//			"posts": blogs,
//		})
//	}
func GetAllPosts(c *gin.Context) {
	page := 1
	limit := 10
	sort := c.DefaultQuery("sort", "desc")
	search := c.Query("search")
	userIDStr := c.Query("user_id")   // optional filter
	includeDraft := c.Query("drafts") // optional: "true" to include drafts

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

	// Only include drafts if requested
	if includeDraft != "true" {
		query = query.Where("published = ?", true)
	}

	// Search filter
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// User filter
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user_id"})
			return
		}
		query = query.Where("user_id = ?", uint(userID))
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to count posts"})
		return
	}

	// Sorting and pagination
	offset := (page - 1) * limit
	order := "created_at desc"
	if sort == "asc" {
		order = "created_at asc"
	}

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
