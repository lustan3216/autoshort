package main

import (
	"autoshort/functions"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// 跑請求前，要先定義的參數
var inputData = InputData{
	Stories: "story1|story2|story3|story3",
	Tags:    "#123, #dsa, #safa",
	Title:   "title",
	Desc:    "description",
}

type InputData struct {
	Stories string `json:"stories"`
	Tags    string `json:"tags"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
}

var templateIDs = map[int]string{
	4:  "295ef471-ff1d-440f-9017-ff328bf756af",
	6:  "868abbd6-9400-432e-a62b-4e25e44c2616",
	9:  "77930739-a006-4732-9407-c2c7e264d61f",
	10: "62cd232f-748f-4612-ad99-4c44202d983f",
	11: "78a1a7e2-e2e0-46d3-8f0c-6ef6ae696e8d",
	12: "af8ab047-ec17-4b23-8610-e9edb915617d",
	13: "c3821f2a-c6dd-4ec8-891d-bef736fe2cc5",
	14: "af8ab047-ec17-4b23-8610-e9edb915617d",
}

const (
	authorization = "Bearer e8601bf029c64990b98967626d07c91de84dc6ac0cdb9812499b0d7df86889e5fadbf247e7a9ba3f763ff39ace832276"
	contentType   = "application/json"
	url           = "https://api.creatomate.com/v1/renders"
)

func main() {
	mp4Url := createVideo()

	for {
		if testVideoURL(mp4Url) {
			uploadToYoutube(mp4Url)
			break // Exit the loop if the video is valid
		}
		time.Sleep(10 * time.Second)
	}
}

func uploadToYoutube(mp4Url string) {
	description := inputData.Tags + "\n\n" + inputData.Desc
	log.Println(mp4Url)
	log.Println(description)
	functions.UploadVideo(
		mp4Url,
		inputData.Title,
		description,
		"27",
		inputData.Tags,
	)
}

func getTemplateID(storyCount int) string {
	if val, ok := templateIDs[storyCount]; ok {
		return val
	}
	log.Fatalf("No template ID found for %d stories", storyCount)
	return ""
}

func createModifications(stories []string) map[string]string {
	modifications := make(map[string]string)
	for index, story := range stories {
		desc := "The scene should have a cinematic feel, without text, vertical, high resolution, eerie tone, a darker color palette, soft lighting, and create a mysterious, manipulative atmosphere."
		modifications[fmt.Sprintf("Image-%d", index+1)] = strings.TrimSpace(story) + ". " + desc
		modifications[fmt.Sprintf("Voiceover-%d", index+1)] = strings.ReplaceAll(strings.TrimSpace(story), ",.:", "")
	}
	return modifications
}

func sendRequest(templateID string, modifications map[string]string) []map[string]interface{} {
	// 定義請求頭
	headers := map[string]string{
		"Authorization": authorization,
		"Content-Type":  contentType,
	}

	// 定義請求體
	body := map[string]interface{}{
		"template_id":   templateID,
		"modifications": modifications,
	}
	bodyJSON, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("Failed to marshal body: %v", err)
	}

	// 發送請求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 解析響應
	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to decode response: %v %v", data, err)
	}

	return data
}

func createVideo() string {
	stories := strings.Split(inputData.Stories, "|")
	storyCount := len(stories)

	templateID := getTemplateID(storyCount)
	modifications := createModifications(stories)
	data := sendRequest(templateID, modifications)

	if len(data) > 0 {
		if url, ok := data[0]["url"].(string); ok {
			return url
		}
	}

	return ""

}

func testVideoURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status: %s\n", resp.Status)
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		fmt.Printf("Expected content type 'video/mp4', but got '%s'\n", contentType)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return false
	}

	if len(body) < 100 {
		// Assuming that a very short response is likely to be an error message
		fmt.Printf("Video is not ready yet, retry after 10 seconds")
		return false
	}

	return true
}
