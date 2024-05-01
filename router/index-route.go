package router

import (
	"learn/httpserver/controller"
	"learn/httpserver/utils"
	"net/url"

	"github.com/gin-gonic/gin"
	// "learn/httpserver/controller"
	hydra "github.com/ory/hydra-client-go/client"
)

var (
	adminURL, _ = url.Parse("http://localhost:4445")
	hydraClient = hydra.NewHTTPClientWithConfig(nil,
		&hydra.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
)

// var userInfo = []repouser.UserInfo{
// 	{
// 		ID:       1,
// 		Email:    "user@example.com",
// 		Password: "password",
// 	},
// 	{
// 		ID:       2,
// 		Email:    "user2@example.com",
// 		Password: "password",
// 	},
// }

func IndexRoute(route *gin.Engine) {

	// gRoute := route.Group("/testGroupPath")

	route.GET("/", controller.Get)

	route.POST("/create", controller.Create)

	route.DELETE("/:id", controller.Delete)

	route.PUT("/:id", controller.Update)

	// route.POST("/login", controller.Login)

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

	// --------- hydra ----------------

	Hydracontroller := controller.Handler{
		HydraAdmin: hydraClient.Admin,
	}

	route.GET("/oauth2/auth", controller.HydraPublicPortCall)
	route.GET("/login", Hydracontroller.AuthGetLogin)
	route.POST("/login", Hydracontroller.AuthPostLogin)
	route.GET("/consent", Hydracontroller.AuthGetConsent)
	route.POST("/consent", Hydracontroller.AuthPostConsent)
	// call hydra token endpoint
	route.POST("/oauth2/token", Hydracontroller.HydraTokenEndpoint)

	// --------- hydra ----------------

}
