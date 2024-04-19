package router

import (
	"learn/httpserver/controller"
	"learn/httpserver/utils"

	"github.com/gin-gonic/gin"
	// "learn/httpserver/controller"
)

func IndexRoute(route *gin.Engine) {

	// gRoute := route.Group("/testGroupPath")

	route.GET("/", controller.Get)

	route.POST("/create", controller.Create)

	route.DELETE("/:id", controller.Delete)

	route.PUT("/:id", controller.Update)

	route.POST("/login", controller.Login)

	// get refresh token and generate new access token
	route.POST("/refresh-token", utils.VerifyToken(1), controller.RefreshToken)

	// 0 for access token, 1 for refresh token
	route.GET("/employee", utils.VerifyToken(0), controller.GetEmployeeData)

	route.POST("/logout", controller.Logout)

	// //to authenticate with jwt
	// route.GET("/auth", middleware.AuthenticateUser, middleware.ValidatePermission, controller.AuthData)

	// route.GET("/session-test", middleware.AuthenticateUser, middleware.ValidatePermission, controller.SessionTest)

	//service
	route.POST("/new-service", controller.AssignNewServiceToUser)

}
