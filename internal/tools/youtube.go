package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// #region Constants/Types

// A slice of valid video types that can be requested from the YouTube API.
var ValidVideoTypes = []string{"video", "short", "vod", "live"}

// The interface for the YouTube service, which is what's used by other modules/files to interact with the YouTube API.
type YTService interface {
	// Depending on the value of `videoType`, retrieves the latest video, short, VOD, or live stream.
	GetLatest(videoType string) (string, error)
}

// The struct that implements the YTService interface.
type ytService struct {
	yt *youtube.Service
}

// #endregion Constants/Types

// #region Helper Functions

// Verify the given videoType parameter is a valid video type by checking if it is in the `ValidVideoTypes` slice.
func ValidateVideoType(videoType string) bool {
	for _, opt := range ValidVideoTypes {
		if opt == videoType {
			return true
		}
	}

	return false
}

// Get the channelName parameter from the request context and validate it. If the channelName parameter is invalid, return an error.
func GetChannelNameParam(c *gin.Context) (string, error) {
	channelName := c.Param("channelName")

	if channelName[0] != '@' {
		return "", fmt.Errorf("invalid channelName parameter, must start with an @ symbol")
	}

	return channelName, nil
}

// Get the videoType parameter from the request context and validate it. If the videoType parameter is invalid, return an error.
func GetVideoTypeParam(c *gin.Context) (string, error) {
	videoType := c.Param("videoType")

	if !ValidateVideoType(videoType) {
		return "", fmt.Errorf("invalid videoType parameter, must be one of: %s", strings.Join(ValidVideoTypes, ","))
	}

	return videoType, nil
}

// #endregion Helper Functions

// #region YTService Implementation/Methods

// Create a new instance of the YTService interface and return it.
func New() YTService {
	ctx := context.Background()
	apiKey := os.Getenv("YT_API_KEY")

	fmt.Println("[tools/youtube#New]: Creating new YTService w/ apiKey = ", apiKey)

	yt, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	return &ytService{yt: yt}
}

// Get the channel ID of the YouTube channel with the given handle.
func (yts *ytService) GetChannelId(handle string) (string, error) {
	res, err := yts.yt.Channels.List([]string{"snippet"}).ForHandle(handle).Do()
	if err != nil {
		log.Fatalf("error getting channel id for channel with \"%v\" handle: %v", handle, err)
	}

	if len(res.Items) > 0 {
		return res.Items[0].Id, nil
	}

	return "", fmt.Errorf("channel not found for handle %s", handle)
}

// Get the latest video, short, VOD, or live stream for the given channel, depending on the value of the `videoType` parameter.
func (yts *ytService) GetLatest(videoType string) (string, error) {
	fmt.Println("[ytService#GetLatest]: videoType = ", videoType)

	yts.yt.Channels.List([]string{"snippet"})

	return "", nil
}

// #endregion YTService Implementation/Methods
