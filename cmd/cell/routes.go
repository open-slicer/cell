package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

		promHandler := promhttp.Handler()
		v2.GET("/metrics", configAuthMiddleware("prometheus.token"), func(c *gin.Context) {
			promHandler.ServeHTTP(c.Writer, c.Request)
		})

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
			locket.PUT("/", configAuthMiddleware("locket.token"), handleLocketPut)
		}
	}
}
