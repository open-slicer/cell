package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRoutes() {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://owo.gg/slicer/cell/-/wikis")
	})

	v2 := r.Group("/api/v2")
	{
		authMiddleware, _ := getAuthMiddleware()
		authBlock := authMiddleware.MiddlewareFunc()

		v2.POST("/users", handleUsersPost)
		v2.GET("/users/:id", authBlock, handleUsersGet)

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}

		locket := v2.Group("/locket")
		{
			locket.Use(locketAuthMiddleware)
			locket.PUT("/", handleLocketPut)
		}
	}
}
