package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Add endpoints to the router, necessary for testing
func SetupRouter(router *gin.Engine, logger *zap.Logger) {
	sugar := logger.Sugar()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World ... ")
	})

	router.POST("/register/user", func(ctx *gin.Context) {
		RegisterUser(ctx, sugar)
	})

	router.POST("/register/workspace", func(ctx *gin.Context) {
		RegisterWorkspace(ctx, sugar)
	})

	router.POST("/register/user_to_workspace", func(ctx *gin.Context) {
		RegisterUserToWorkspace(ctx, sugar)
	})

	// TODO: [ ] Send
	router.POST("/update/me", func(ctx *gin.Context) {
		UserIPCheck(ctx, sugar)
	})

	// Returns all workspace owner's IP - (In Future allow all the IPs of people that have latest workspace)
	// TODO: [ ] Test this code :)
	// TODO: [ ] A way to load dummy data for test in test db
	router.GET("/whats/new", func(ctx *gin.Context) {
		GetAllMyConnectedWorkspaceInfo(ctx, sugar)
	})
}
