package handler

import (
	"math/rand"
	"strconv"

	"go.uber.org/zap"

	"github.com/ButterHost69/PKr-server/db"
	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context, sugar *zap.SugaredLogger) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	tagId := rand.Intn(9000) + 1000
	username = username + "#" + strconv.Itoa(tagId)

	if err := db.CreateNewUser(username, password); err != nil {
		sugar.Error(err)
		ctx.JSON(500, gin.H{
			"response": "internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"response": "success",
		"username": username,
	})
}

// TODO: [ ] Check if user already has a workspace of that name already. Return with another type of response if so eg. Workspace Exists
func RegisterWorkspace(ctx *gin.Context, sugar *zap.SugaredLogger) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	workspace_name := ctx.PostForm("workspace_name")
	
	if username == "" || password == "" || workspace_name == ""{
		ctx.JSON(203, gin.H{
			"response": "incorrect parameters",
		})	
	}

	ok, err := db.RegisterNewWorkspace(username, password, workspace_name)
	if err != nil {
		sugar.Error(err)
		ctx.JSON(500, gin.H{
			"response": "internal server error",
		})	
		return
	}

	if ok {
		ctx.JSON(200, gin.H{
			"response": "success",
		})	
		return
	}

	ctx.JSON(201, gin.H{
		"response": "authentication error",
	})	
}
