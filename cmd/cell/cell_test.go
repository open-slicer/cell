package main

import (
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"net/http"
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

func getExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(r),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
