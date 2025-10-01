package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	docs "github.com/raihaninkam/finalPhase3/docs"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	// router.Use(middlewares.CORSMiddleware)

	InitAuthRouter(router, db, rdb)

	InitPostRouter(router, db, rdb)

	InitFollowsRouter(router, db, rdb)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Rute Salah",
			"Status":  "Rute Tidak Ditemukan",
		})
	})
	return router

}
