package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/raihaninkam/finalPhase3/internals/utils"
)

type CommentHandler struct {
	cr *repositories.CommentRepository
}

func NewCommentHandler(cr *repositories.CommentRepository) *CommentHandler {
	return &CommentHandler{cr: cr}
}

func (ch *CommentHandler) CreateComment(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Post ID is required",
		})
		return
	}

	var req models.CommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	comment := &models.Comment{
		UserId:  userID,
		PostId:  postID,
		Content: req.Content,
	}

	newComment, err := ch.cr.CreateComment(ctx, comment)
	if err != nil {
		log.Println("Error creating comment:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    newComment,
	})
}

func (ch *CommentHandler) GetPostComments(ctx *gin.Context) {
	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Post ID is required",
		})
		return
	}

	comments, err := ch.cr.GetPostComments(ctx, postID)
	if err != nil {
		log.Println("Error getting comments:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
	})
}

func (ch *CommentHandler) DeleteComment(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	commentID := ctx.Param("id")
	if commentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Comment ID is required",
		})
		return
	}

	if err := ch.cr.DeleteComment(ctx, commentID, userID); err != nil {
		log.Println("Error deleting comment:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment deleted successfully",
	})
}
