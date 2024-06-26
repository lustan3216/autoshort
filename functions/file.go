package functions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// downloadFile downloads a file from the given URL and saves it to the given file path.
func DownloadFile(filePath string, url string) error {
	// Create the file
	log.Printf("Downloading file from %s ......\n", url)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
