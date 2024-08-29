package api

import (
	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
)

func Serve() {
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	api := r.Group("/api")
	api.Use(AuthMiddleware())
	{
		api.GET("/hello", helloHandler)
		api.GET("/group", groupHandler)
		api.PATCH("/group", groupScore)
		api.GET("/anime", animesHandler)
		api.PATCH("/anime/:id", blackAnime)

		r.POST("/api/login", loginHandler)

	}

	if err := r.Run(":8098"); err != nil {
		panic(err)
	}

}
