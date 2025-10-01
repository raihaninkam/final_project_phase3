package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/handlers"
	middleware "github.com/raihaninkam/finalPhase3/internals/middlewares"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitFollowsRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	followRouter := router.Group("")
	followRepository := repositories.NewFollowsRepository(db, rdb)
	followHandler := handlers.NewFollowHandler(followRepository)

	followRouter.GET("/following", middleware.VerifyToken(rdb), followHandler.GetFollowing)
	followRouter.POST("/follow/:id", middleware.VerifyToken(rdb), followHandler.Follow)

}
