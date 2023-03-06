// run code locally with  CompileDaemon -command="./go-jwt"
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/viru56/go-jwt/controllers"
	"github.com/viru56/go-jwt/intializers"
	"github.com/viru56/go-jwt/middleware"
)

func init() {
	intializers.LoadEnvVariables()
	intializers.ConnectToDb()
	intializers.SyncDatabase()
}

func main () {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/user", middleware.RequireAuth, controllers.GetUser)
	r.POST("/logout", middleware.RequireAuth, controllers.Logout)
	r.Run()
}