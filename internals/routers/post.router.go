package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/handlers"
	middleware "github.com/raihaninkam/finalPhase3/internals/middlewares"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitPostRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	postRouter := router.Group("")
	postRepository := repositories.NewPostRepository(db, rdb)
	postHandler := handlers.NewPostHandler(postRepository)

	// posting
	postRouter.POST("/post", middleware.VerifyToken(rdb), postHandler.CreatePost)
	postRouter.GET("/post", middleware.VerifyToken(rdb), postHandler.GetFeed)

	// like
	likeRepository := repositories.NewLikeRepository(db, rdb)
	likeHandler := handlers.NewLikeHandler(likeRepository)
	postRouter.POST("/post/:id/like", middleware.VerifyToken(rdb), likeHandler.LikePost)
	postRouter.DELETE("/post/:id/like", middleware.VerifyToken(rdb), likeHandler.UnlikePost)

	// comment
	commentRepository := repositories.NewCommentRepository(db, rdb)
	commentHandler := handlers.NewCommentHandler(commentRepository)
	postRouter.GET("/post/:id/comment", middleware.VerifyToken(rdb), commentHandler.GetPostComments)
	postRouter.POST("/post/:id/comment", middleware.VerifyToken(rdb), commentHandler.CreateComment)
	postRouter.DELETE("/comment/:id", middleware.VerifyToken(rdb), commentHandler.DeleteComment)

}
