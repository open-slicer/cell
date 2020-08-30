package main

import (
	"context"
	"net/http"
	"testing"
)

const (
	username = "test"
	password = "_--fgsdLjhKf--_"
)

func TestUsers(t *testing.T) {
	e := getExpect(t)

	// Delete the user if it already exists. This is to avoid a conflict status code.
	_, err := pg.Exec(context.Background(), "DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Fatalf("%e", err)
	}

	id := e.POST("/api/v2/users").WithJSON(userInsertion{
		Username:    username,
		DisplayName: "Test User",
		PublicKey:   []byte("..."),
		Password:    password,
	}).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		Value("data").Object().
		Value("id").String().Raw()

	token := "Bearer " + e.POST("/api/v2/auth/login").WithJSON(userLogin{
		Username: username,
		Password: password,
	}).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("token").String().Raw()

	token = "Bearer " + e.GET("/api/v2/auth/refresh").WithHeader("Authorization", token).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("token").String().Raw()

	e.GET("/api/v2/users/"+id).WithHeader("Authorization", token).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		ValueEqual("id", id)
}
