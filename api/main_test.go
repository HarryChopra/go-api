package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Cleaner test output
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
