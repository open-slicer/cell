package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// TODO: Make this stricter.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://slicer.softwares.software/cell")
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
			users.POST("", handleUsersPost)
			users.GET("/:id", authBlock, handleUsersGet)
		}

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}

		lockets := v2.Group("/lockets")
		{
			lockets.GET("", authBlock, handleLocketGet)
			lockets.PUT("", configAuthMiddleware("locket.token"), handleLocketPut)
		}

		channels := v2.Group("/channels")
		{
			channels.Use(authBlock)
			channels.POST("", handleChannelsPOST)
		}
	}

	return r
}
