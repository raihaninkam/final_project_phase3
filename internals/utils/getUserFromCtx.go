package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/raihaninkam/finalPhase3/pkg"
)

func GetUserFromCtx(ctx *gin.Context) (string, error) {
	claims, ok := ctx.Get("claims")
	if !ok {
		return "", errors.New("claims not found in context, token might be missing")
	}

	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return "", errors.New("invalid claims format")
	}

	return userClaims.UserId, nil
}
