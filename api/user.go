package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"ws-home-backend/business"
)

func GetUserInfoById(ctx *gin.Context) {

	value := ctx.Query("userId")
	userId, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Get user info by id" + fmt.Sprintf("%d", userId),
		"data": business.GetUserById(int32(userId)),
	})
}
