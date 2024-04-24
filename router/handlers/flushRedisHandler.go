package handlers

import (
	"Institution/redis"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FlushRedisHandler(ctx *gin.Context) {
	redisClient := redis.GetClient()
	redisClient.FlushAll(context.Background())
	ctx.JSON(http.StatusOK, gin.H{})
}
