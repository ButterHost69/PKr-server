package handler

import (
	"math/rand"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/ButterHost69/PKr-server/db"
	"github.com/gin-gonic/gin"
)

func getIpAddress(ctx *gin.Context) (string, string) {
	remoteAddr := ctx.Request.RemoteAddr
	address := strings.Split(remoteAddr, ":")

	return address[0], address[1]
}

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

	ip_addr, port := getIpAddress(ctx)
	if err := db.UpdateUserIP(username, password, ip_addr, port); err != nil {
		ctx.JSON(203, gin.H{
			"response": "could not update user ip",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"response": "success",
		"username": username,
	})
}

// TODO: [X] Check if user already has a workspace of that name already. Return with another type of response if so eg. Workspace Exists
// TODO: [ ] Check if the above todo works, no time for today :)
func RegisterWorkspace(ctx *gin.Context, sugar *zap.SugaredLogger) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	workspace_name := ctx.PostForm("workspace_name")

	if username == "" || password == "" || workspace_name == "" {
		ctx.JSON(203, gin.H{
			"response": "incorrect parameters",
		})
	}

	workspaceList, err := db.GetWorkspaceList(username)
	if err != nil {
		sugar.Error(err)
		ctx.JSON(500, gin.H{
			"response": "internal server error",
		})
		return
	}

	for _, val := range workspaceList {
		if val == workspace_name {
			ctx.JSON(200, gin.H{
				"response": "workspace already exist",
			})
			return
		}
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

// TODO : [ ] Check if the Below code Works, havent tested it :)
// TODO : [ ] Check if the IP of the Requester matches in the LastUserIP Database
func RegisterUserToWorkspace(ctx *gin.Context, sugar *zap.SugaredLogger) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	workspace_name := ctx.PostForm("workspace_name")
	connection_username := ctx.PostForm("connection_username")

	if username == "" || password == "" || workspace_name == "" || connection_username == "" {
		ctx.JSON(203, gin.H{
			"response": "incorrect parameters",
		})
		return
	}

	ok, err := db.RegisterUserToWorkspace(username, password, workspace_name, connection_username)
	
	switch ok {
	case 0:
		ctx.JSON(200, gin.H{
			"response": "success",
		})
		return
	case 1:
		ctx.JSON(201, gin.H{
			"response": "authentication error",
		})
		return
	case 2:
		ctx.JSON(201, gin.H{
			"response": "workspace doesn't exist",
		})
		return
	case 5:
		sugar.Error(err)
		ctx.JSON(201, gin.H{
			"response": "internal server error",
		})
		return
	}

}

// TODO: [ ] Make so that user has to prove username password only once do something like session token...
func UserIPCheck(ctx *gin.Context, sugar *zap.SugaredLogger) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	port := ctx.PostForm("port")

	var ipaddr string
	if username == "" {
		ctx.JSON(203, gin.H{
			"response": "incorrect parameters",
		})
		return
	}

	if port == "" {
		ipaddr, port = getIpAddress(ctx)
	} else {
		ipaddr, _ = getIpAddress(ctx)
	}

	if err := db.UpdateUserIP(username, password, ipaddr, port); err != nil {
		ctx.JSON(203, gin.H{
			"response": "could not update user ip",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"response": "success",
	})
}
