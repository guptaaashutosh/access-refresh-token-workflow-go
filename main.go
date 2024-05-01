package main

import (
	"learn/httpserver/router"
	"learn/httpserver/setup"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	r := setupRouter()

	r.Use(CORSMiddleware())

	//Register the standard HandlerFuncs from the net/http/pprof package with the provided gin.Engine.
	pprof.Register(r)
	// pprof.Register(router, &pprof.Options{
	// 	// default is "debug/pprof"
	// 	RoutePrefix: "debug/pprof",
	// })

	r.Run()

}
