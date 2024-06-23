package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dylanmei/iso8601"
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
	GetLatestItem(channelId string, videoType string) (*youtube.Video, error)

	// Retrieves the channel ID of the YouTube channel with the given handle.
	GetChannelId(handle string) (string, error)

	// Retrieves the full video object for the video with the given ID.
	GetVideoById(videoId string) (*youtube.Video, error)
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

// Determines if the given YouTube item is an official "YouTube Short". The API doesn't make this distinction, so we have to infer it based on a few criteria.
func isVideoShort(item youtube.Video) bool {
	title := strings.ToLower(item.Snippet.Title)

	// fmt.Printf("[isVideoShort]: title = %s\n", title)

	if strings.Contains(title, "#short") {
		return true
	}

	// fmt.Println("[isVideoShort]: item.ContentDetails.Duration = ", item.ContentDetails.Duration)
	duration, err := iso8601.ParseDuration(item.ContentDetails.Duration)
	if err != nil {
		log.Fatalf("[isVideoShort]: error parsing duration: %v", err)
	}

	if duration.Seconds() < 60 && duration.Seconds() > 0 {
		return true
	}

	// fmt.Printf("[isVideoShort]: duration = %v\nduration.Seconds() = %v\n", duration, duration.Seconds())
	// fmt.Printf("[isVideoShort]: fancyDuration = %v\n", fancyDuration)

	return false
}

// #endregion Helper Functions

// #region YTService Implementation/Methods

// Create a new instance of the YTService interface and return it.
func NewYTService() YTService {
	ctx := context.Background()
	apiKey := os.Getenv("YT_API_KEY")

	fmt.Println("[tools/youtube#New]: Creating new YTService w/ apiKey = ", apiKey)

	yt, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	return &ytService{yt: yt}
}

func (yts *ytService) GetVideoById(videoId string) (*youtube.Video, error) {
	res, err := yts.yt.Videos.List([]string{"snippet", "contentDetails", "liveStreamingDetails"}).Id(videoId).Do()
	if err != nil {
		log.Fatalf("error getting video with \"%v\" id: %v", videoId, err)
	}

	if len(res.Items) > 0 {
		return res.Items[0], nil
	}

	return nil, fmt.Errorf("video not found for id %s", videoId)
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

// Get the uploads playlist id for the channel with the given id.
func (yts *ytService) GetUploadPlaylistId(channelId string) (string, error) {
	res, err := yts.yt.Channels.List([]string{"contentDetails"}).Id(channelId).Do()
	if err != nil {
		log.Fatalf("error getting uploads playlist id for channel with \"%v\" id: %v", channelId, err)
	}

	if len(res.Items) > 0 {
		return res.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
	}

	return "", fmt.Errorf("uploads playlist not found for channel with id %s", channelId)
}

// Get the latest video, short, VOD, or live stream for the given channel, depending on the value of the `videoType` parameter.
func (yts *ytService) GetLatestItem(channelId string, videoType string) (*youtube.Video, error) {
	uploadsPlaylistId, err := yts.GetUploadPlaylistId(channelId)
	if err != nil {
		return nil, err
	}

	items, err := yts.yt.PlaylistItems.List([]string{"snippet"}).PlaylistId(uploadsPlaylistId).MaxResults(50).Do()
	if err != nil {
		log.Fatalf("error getting playlist items for channel with \"%v\" id: %v", channelId, err)
		return nil, err
	}

	for _, item := range items.Items {
		video, err := yts.GetVideoById(item.Snippet.ResourceId.VideoId)
		if err != nil {
			fmt.Printf("error getting video with \"%v\" id: %v", item.Snippet.ResourceId.VideoId, err)
		}
		//  else {
		// 	// fmt.Printf("[GetLatest]: Successfully retrieved full video object for video with id \"%v\"\n", videoId)
		// }

		// fmt.Printf("[GetLatest]: video.LiveStreamingDetails.ActualStartTime = %s\n", video.LiveStreamingDetails.ActualStartTime)
		// fmt.Printf("[GetLatest]: video.LiveStreamingDetails.ActualEndTime = %s\n", video.LiveStreamingDetails.ActualEndTime)
		// fmt.Printf("[ytService#GetLatest]: i = %v & isShort = %v\n", i, isShort)

		// If the videoType is short, first check if the video is a short video.
		if videoType == "short" && isVideoShort(*video) {
			return video, nil
		} else if videoType == "live" && video.Snippet.LiveBroadcastContent == "upcoming" {
			return video, nil
			// If the videoType is vod, check if the video.liveStreamingDetails.actualStartTime and video.liveStreamingDetails.actualEndTime are both not empty.
		} else if videoType == "vod" && video.LiveStreamingDetails.ActualStartTime != "" && video.LiveStreamingDetails.ActualEndTime != "" {
			return video, nil
		} else {
			return video, nil
		}
	}

	return nil, fmt.Errorf("no video found of type %s for channel with id %s", videoType, channelId)
}

// #endregion YTService Implementation/Methods
