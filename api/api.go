package api

import (
	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
)

func Serve() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	r.GET("/hello", helloHandler)
	r.POST("/api/login", loginHandler)
	r.GET("/api/verify", verifyHandler)
	api := r.Group("/api")

	api.Use(AuthMiddleware())
	{
		api.GET("/group", groupHandler)
		api.PATCH("/group", groupScore)
		api.GET("/anime", animesHandler)
		api.PATCH("/anime/:id", blackAnime)
	}

	if err := r.Run(":8099"); err != nil {
		panic(err)
	}
}
