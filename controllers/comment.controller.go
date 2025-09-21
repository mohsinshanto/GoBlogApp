package controllers

import (
	"BlogApp/config"
	"BlogApp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	ID       uint   `json:"id"`
	Content  string `json:"content"`
	PostID   uint   `json:"post_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// authenticated user
func CreateComment(c *gin.Context) {
	// Get user_id from context (set in middleware)
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Safely assert user_id to float64 and convert to uint
	floatID, ok := uid.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	userID := uint(floatID)

	// Input struct with validation
	var input struct {
		Content string `json:"content" binding:"required"`
		PostID  uint   `json:"post_id" binding:"required"`
	}

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the post exists
	var post models.Blog
	if err := config.DB.First(&post, input.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Create the comment
	comment := models.Comment{
		Content: input.Content,
		PostID:  input.PostID,
		UserID:  userID,
	}

	// Save the comment
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusCreated, comment)
}

func UpdateComment(c *gin.Context) {
	// Get user_id from context
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Safely assert user_id from float64 to uint
	floatID, ok := uid.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	userID := uint(floatID)

	// Parse the comment ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment ID"})
		return
	}

	// Fetch comment from DB
	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if the logged-in user owns the comment
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own comment"})
		return
	}

	// Bind new content
	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update content
	comment.Content = input.Content
	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	// Get user_id from context
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Safely assert user_id to float64, then convert to uint
	floatID, ok := uid.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	userID := uint(floatID)

	// Parse the comment ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment ID"})
		return
	}
	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comment"})
		return
	}

	// Delete comment
	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// Any users
func GetComment(c *gin.Context) {
	id := c.Param("id")

	// Convert id to uint
	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, uint(commentID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func GetAllComments(c *gin.Context) {
	var comments []models.Comment
	if err := config.DB.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
func GetCommentCount(c *gin.Context) {
	postIDParam := c.Param("post_id")

	// Convert postID from string to uint
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var count int64
	if err := config.DB.Model(&models.Comment{}).
		Where("post_id = ?", uint(postID)).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id":        postID,
		"total_comments": count,
	})
}

// Get all comments for a specific post
func GetCommentsByPost(c *gin.Context) {
	postIDParam := c.Param("post_id")

	var comments []models.Comment
	if err := config.DB.
		Preload("User").
		Where("post_id = ?", postIDParam).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map DB model to response model
	var response []CommentResponse
	for _, comment := range comments {
		response = append(response, CommentResponse{
			ID:       comment.ID,
			Content:  comment.Content,
			PostID:   comment.PostID,
			UserID:   comment.UserID,
			Username: comment.User.Username,
		})
	}

	c.JSON(http.StatusOK, response)
}
