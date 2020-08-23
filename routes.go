package main

func setupRoutes() {
	authMiddleware, _ := getAuthMiddleware()
	authBlock := authMiddleware.MiddlewareFunc()

	v2 := r.Group("/api/v2")
	{
		v2.POST("/users", handleUsersPost)
		v2.GET("/users/:id", authBlock, handleUsersGet)

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}
	}
}
