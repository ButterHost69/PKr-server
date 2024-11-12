package handler

import (
	"math/rand"
	"strconv"

	"github.com/ButterHost69/PKr-server/db"
	"github.com/gin-gonic/gin"
)

func Register_User(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	tagId := rand.Intn(9000) + 1000
	username = username + "#" + strconv.Itoa(tagId)

	if err := db.CreateNewUser(username, password); err != nil {
		ctx.JSON(500, gin.H{
			"response":"internal server error",
		})
		return
	}


	ctx.JSON(200, gin.H{
		"response":"success",
		"username": username,
	})
}