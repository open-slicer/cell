package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	readConfig()
	dbConnect()
	r = setupRouter()
}
