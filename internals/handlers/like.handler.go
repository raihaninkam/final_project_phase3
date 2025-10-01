package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/raihaninkam/finalPhase3/internals/utils"
)

type LikeHandler struct {
	lr *repositories.LikeRepository
}

func NewLikeHandler(lr *repositories.LikeRepository) *LikeHandler {
	return &LikeHandler{lr: lr}
}

func (lh *LikeHandler) LikePost(ctx *gin.Context) {
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

	if err := lh.lr.LikePost(ctx, userID, postID); err != nil {
		log.Println("Error liking post:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post liked successfully",
	})
}

func (lh *LikeHandler) UnlikePost(ctx *gin.Context) {
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

	if err := lh.lr.UnlikePost(ctx, userID, postID); err != nil {
		log.Println("Error unliking post:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post unliked successfully",
	})
}
