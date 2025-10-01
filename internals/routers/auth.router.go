package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raihaninkam/finalPhase3/internals/handlers"
	middleware "github.com/raihaninkam/finalPhase3/internals/middlewares"
	"github.com/raihaninkam/finalPhase3/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db, rdb)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/login", authHandler.Login)

	authRouter.PATCH("/profile", middleware.VerifyToken(rdb), authHandler.UpdateProfile)

}
