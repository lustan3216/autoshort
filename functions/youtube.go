// Sample Go code for user authorization

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"
)

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func UploadVideo(videoUrl string, title string, description string, category string, keywords string) {
	tmpFilePath := "video.mp4"

	client := getClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")

	err = DownloadFile(tmpFilePath, videoUrl)
	handleError(err, "")

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "public",
			MadeForKids:   false,
		},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(tmpFilePath)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening %v: %v", tmpFilePath, err)
	}

	response, err := call.Media(file).Do()
	handleError(err, "")

	DeleteFile(tmpFilePath)
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)

}

func main() {
	UploadVideo(
		"https://cdn.creatomate.com/renders/1179d3de-4a35-4c58-b469-11d54c97e5b3.mp4",
		"My Video",
		"This is a test video",
		"22",
		"test, video, upload",
	)
}
