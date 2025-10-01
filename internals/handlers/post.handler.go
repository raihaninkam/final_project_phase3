package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/internals/models"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/raihaninkam/finalPhase3/internals/utils"
)

type PostHandler struct {
	pr *repositories.PostRepository
}

func NewPostHandler(pr *repositories.PostRepository) *PostHandler {
	return &PostHandler{pr: pr}
}

func (ph *PostHandler) CreatePost(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Ambil form field text
	content := ctx.PostForm("content_text")

	// Ambil file
	file, err := ctx.FormFile("image")
	var imageUrl string
	if err == nil {
		uploadedFile, err := utils.FileUpload(ctx, file, "post")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		imageUrl = "/public/" + uploadedFile
	}

	// Validasi: minimal harus ada content atau image
	if content == "" && imageUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Content or image is required",
		})
		return
	}

	// Buat object Posts
	post := &models.Posts{
		UserId:   userID,
		Content:  content,
		ImageUrl: imageUrl,
	}

	newPost, err := ph.pr.CreatePost(ctx, post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    newPost,
	})
}

func (ph *PostHandler) GetFeed(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	// Get pagination params
	limit := 20
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	posts, err := ph.pr.GetFollowingPosts(ctx, userID, limit, offset)
	if err != nil {
		log.Println("Error getting feed:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
	})
}
