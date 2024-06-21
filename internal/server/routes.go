package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/4lch4/YouToob/internal/templates"
	"github.com/4lch4/YouToob/internal/tools"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func getRootPathHelp(baseUrl string) (string, error) {
	tmpl, err := template.New("BaseRouteHelp").Parse(templates.GetBaseRouteHelp())
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, &templates.BaseRouteHelp{BaseURL: baseUrl, ChannelName: "sneeziu"})
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buff.String(), nil
}

// Build and return the info string used for the root path, which provides information about the available endpoints.
// func getRootPathHelp(urlBase string) string {
// 	return fmt.Sprintf(`
// This is an informational endpoint. The actual endpoints that are available are as follows:

// - GET %s/:channelName
// 	- Returns a JSON object where each key is a valid videoType and the value is the URL for that videoType.
// 	- Example: GET %s/:channelName
// - GET %s/:channelName/:videoType
// 	- Returns a string with the title and URL of the latest video of the specified type.
// 	- Example response: She Did WHAT To Her Twitch Chat?! #shorts - https://youtu.be/gge-pBrHBzo
//     `, urlBase, urlBase, urlBase,
// 	)
// }

// Register the routes for the server to their respective handlers.
func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.GET("/", s.handleBaseRoute)

	// routesGroup := router.Group(os.Getenv("API_BASE_PATH"))

	// router.GET("/", s.helloWorldHandler)
	router.GET("/:channelName", s.handleNoVideoType)
	router.GET("/:channelName/:videoType", s.handleLatestContent)

	return router
}

//#region Handlers

func (s *Server) handleBaseRoute(ctx *gin.Context) {
	helpData, err := getRootPathHelp(tools.GetBaseUrl(ctx))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error generating help text: \n\n%s\n", err)
		return
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, helpData)
}

// Handle requests that do not specify a video type. Return the available video types for the given channel.
func (s *Server) handleNoVideoType(c *gin.Context) {
	channelName, err := tools.GetChannelNameParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseUrl := tools.GetBaseUrl(c)

	c.JSON(http.StatusOK, gin.H{
		"vod":   fmt.Sprintf("%s/%s/vod", baseUrl, channelName),
		"live":  fmt.Sprintf("%s/%s/live", baseUrl, channelName),
		"short": fmt.Sprintf("%s/%s/short", baseUrl, channelName),
		"video": fmt.Sprintf("%s/%s/video", baseUrl, channelName),
	})
}

func (s *Server) handleLatestContent(c *gin.Context) {
	channelName, err := tools.GetChannelNameParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	videoType, err := tools.GetVideoTypeParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	fmt.Printf("[routes#GetLatest]: videoType of \"%s\" & channelName of \"%s\" are valid\n", videoType, channelName)

	c.JSON(http.StatusOK, gin.H{
		"channelName": channelName,
		"videoType":   videoType,
	})
}

func (s *Server) helloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

//#endregion Handlers
