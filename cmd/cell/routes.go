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

		users := v2.Group("/users")
		{
			users.POST("/", handleUsersPost)
			users.GET("/:id", authBlock, handleUsersGet)
		}

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}

		locket := v2.Group("/lockets")
		{
			locket.GET("/", handleLocketGet)
			locket.PUT("/", locketAuthMiddleware, handleLocketPut)
		}
	}
}
