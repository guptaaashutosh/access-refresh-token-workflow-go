package main

import (
	"learn/httpserver/controller"
	"learn/httpserver/middleware"
	"learn/httpserver/router"
	"learn/httpserver/setup"
	"learn/httpserver/utils"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// func init() {
// 	setup.LoadEnvVariable()
// }

func setupRouter() *gin.Engine {

	setup.LoadEnvVariable()

	r := gin.Default()
	redisClient, _ := utils.ConnectRedis()

	r.Use(middleware.SetRedisClientToContext(redisClient))

	//on server start it put all the data in redis
	go controller.PutAllDataInRedis(redisClient)

	store := cookie.NewStore([]byte(os.Getenv("SESSION_KEY")))
	r.Use(sessions.Sessions("my-session", store))

	router.IndexRoute(r)

	return r
}

func main() {

	r := setupRouter()

	r.Run() 

}
