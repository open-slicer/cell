package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"testing"
)

func TestUsers(t *testing.T) {
	e := getExpect(t)

	// Delete the user if it already exists. This is to avoid a conflict status code.
	_, err := mng.users.DeleteOne(context.Background(), bson.M{
		"username": "test",
	})
	if err != nil {
		t.Fatalf("%e", err)
	}

	username := "test"
	password := "_--fgsdLjhKf--_"
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

	token := e.POST("/api/v2/auth/login").WithJSON(userLogin{
		Username: username,
		Password: password,
	}).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("token").String().Raw()

	e.GET("/api/v2/users/"+id).WithHeader("Authorization", "Bearer "+token).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		ValueEqual("id", id)
}
