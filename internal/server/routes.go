package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/4lch4/YouToob/internal/templates"
	"github.com/4lch4/YouToob/internal/tools"
	youtube "github.com/4lch4/YouToob/internal/tools/yt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var yt = youtube.NewYTService()

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

// Register the routes for the server to their respective handlers.
func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.GET("/", s.handleBaseRoute)
	router.GET("/:channelName", s.handleNoVideoType)
	router.GET("/:channelName/:videoType", s.handleLatestContent)

	return router
}

//#region Route Handlers

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
	channelName, err := youtube.GetChannelNameParam(c)
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
	channelName, err := youtube.GetChannelNameParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	videoType, err := youtube.GetVideoTypeParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	fmt.Printf("[routes#GetLatest]: videoType of \"%s\" & channelName of \"%s\" are valid\n", videoType, channelName)

	channel, err := yt.GetChannelId(channelName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("[routes#GetLatest]: Channel ID is ", channel)

	latest, err := yt.GetLatestItem(channel, videoType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channelName": channelName,
		"videoType":   videoType,
		"videoId":     latest,
	})
}

//#endregion Route Handlers
