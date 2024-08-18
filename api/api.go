package api

import (
	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
)

func Serve() {
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	api := r.Group("/api")
	{
		api.GET("/hello", helloHandler)
		api.GET("/group", groupHandler)
		api.GET("/anime", animesHandler)
		api.PATCH("/anime/:id", blackAnime)
	}

	if err := r.Run(":8099"); err != nil {
		panic(err)
	}

}
