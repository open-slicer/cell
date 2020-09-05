package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// TODO: Make this stricter.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetStringSlice("cors.allowed_origins"),
		AllowMethods:     viper.GetStringSlice("cors.alloewd_methods"),
		AllowHeaders:     viper.GetStringSlice("cors.allowed_headers"),
		ExposeHeaders:    viper.GetStringSlice("cors.exposed_headers"),
		AllowCredentials: viper.GetBool("cors.allow_credentials"),
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
			users.POST("", handleUsersPOST)
			specific := users.Group("/:user")
			{
				specific.GET("", authBlock, handleUsersGET)
			}
		}

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}

		lockets := v2.Group("/lockets")
		{
			lockets.GET("", authBlock, handleLocketsGET)
			lockets.PUT("", configAuthMiddleware("locket.token"), handleLocketsPUT)
		}

		invites := v2.Group("/invites")
		{
			invites.Use(authBlock)
			specific := invites.Group("/:invite")
			{
				specific.GET("", handleInvitesGET)
				specific.GET("/accept", handleInvitesAcceptGET)
			}
		}

		channels := v2.Group("/channels")
		{
			channels.Use(authBlock)
			channels.POST("", handleChannelsPOST)

			specific := channels.Group("/:channel")
			{
				specific.GET("", handleChannelsGET)
				specific.POST("/invites", handleInvitesPOST)
				members := specific.Group("/members")
				{
					members.GET("/:member", handleMembersGET)
				}
			}
		}
	}

	return r
}
