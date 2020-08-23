package main

func setupRoutes() {
	authMiddleware, _ := getAuthMiddleware()

	v2 := r.Group("/api/v2")
	{
		v2.POST("/users", handleUsersPost)

		auth := v2.Group("/auth")
		{
			auth.POST("/login", authMiddleware.LoginHandler)
			auth.GET("/refresh", authMiddleware.RefreshHandler)
		}
	}
}
