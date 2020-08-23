package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	}
}
