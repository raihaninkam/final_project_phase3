package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/raihaninkam/finalPhase3/internals/utils"
)

type FollowHandler struct {
	fr *repositories.FollowRepository
}

func NewFollowHandler(fr *repositories.FollowRepository) *FollowHandler {
	return &FollowHandler{fr: fr}
}

func (fh *FollowHandler) GetFollowing(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	users, err := fh.fr.GetFollowing(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

func (fh *FollowHandler) Follow(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	followingID := ctx.Param("id")
	if followingID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	if userID == followingID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Cannot follow yourself",
		})
		return
	}

	if err := fh.fr.Follow(ctx, userID, followingID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to follow user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully followed user",
	})
}
