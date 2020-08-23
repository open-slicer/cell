package main

func init() {
	v2 := r.Group("/api/v2")
	{
		v2.POST("/users", handleUsersPost)
	}
}
