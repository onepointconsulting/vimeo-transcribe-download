package vimeo_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	model "vimeo-transcriber-model"
)

func executeRequest(config *model.Config, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", config.VimeoPersonalAccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return result, nil
}

func GetUser(config *model.Config) (*model.User, error) {
	result, err := executeRequest(config, "https://api.vimeo.com/me")
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &model.User{
		ID:       result["uri"].(string),
		Username: result["name"].(string),
	}, nil
}

func GetTextTracks(config *model.Config, videoId string) ([]string, error) {
	result, err := executeRequest(config, fmt.Sprintf("https://api.vimeo.com/videos/%s/texttracks", videoId))

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	tracks := result["data"].([]interface{})
	links := []string{}
	for _, element := range tracks {
		elementMap := element.(map[string]interface{})
		if hlsLink, ok := elementMap["hls_link"].(string); ok {
			links = append(links, hlsLink)
		} else {
			log.Printf("hls_link not found or not a string")
		}
	}
	return links, nil
}

func GetVttFile(url string, fileName string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting vtt file: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	// Create the directory if it doesn't exist
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}
	os.WriteFile(fileName, bodyBytes, 0644)
	return fileName, nil
}
