package main

import (
	"learn/httpserver/router"
	"learn/httpserver/setup"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/pprof"
)

// func init() {
// 	setup.LoadEnvVariable()
// }

func setupRouter() *gin.Engine {

	setup.LoadEnvVariable()

	r := gin.Default()
	// redisClient, _ := utils.ConnectRedis()

	// r.Use(middleware.SetRedisClientToContext(redisClient))

	// //on server start it put all the data in redis
	// go controller.PutAllDataInRedis(redisClient)

	// store := cookie.NewStore([]byte(os.Getenv("SESSION_KEY")))
	// r.Use(sessions.Sessions("my-session", store))

	router.IndexRoute(r)

	return r
}

func main() {

	r := setupRouter()

	//Register the standard HandlerFuncs from the net/http/pprof package with the provided gin.Engine.
	pprof.Register(r)
	// pprof.Register(router, &pprof.Options{
	// 	// default is "debug/pprof"
	// 	RoutePrefix: "debug/pprof",
	// })

	r.Run()

}
