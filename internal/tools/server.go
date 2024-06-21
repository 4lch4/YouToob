package tools

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Get the base URL for the current request. For example:
//
// - http://localhost:8080/@renniesaurus/live becomes http://localhost:8080
//
// - https://yt.4lch4.com/@renniesaurus/live becomes https://yt.4lch4.com
func GetBaseUrl(c *gin.Context) string {
	var prefix string

	if c.Request.TLS == nil {
		prefix = "http"
	} else {
		prefix = "https"
	}

	return fmt.Sprintf("%s://%s", prefix, c.Request.Host)
}
