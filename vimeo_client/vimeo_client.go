package vimeo_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("Too many requests. Sleeping for %d seconds.", config.RetryDelay)
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
			return executeRequest(config, url)
		}
		if resp.StatusCode == http.StatusNotFound {
			return map[string]interface{}{"error": "not found"}, nil
		}
		return nil, fmt.Errorf("unexpected status code: %d for %s", resp.StatusCode, url)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	// log.Println(string(bodyBytes))

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

func DumpUser(config *model.Config) (string, error) {
	result, err := executeRequest(config, "https://api.vimeo.com/me")
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}
	return dumpJson(result)
}

func GetTextTracks(config *model.Config, videoId string) ([]string, error) {
	result, err := executeRequest(config, fmt.Sprintf("https://api.vimeo.com/videos/%s/texttracks", videoId))

	if err != nil {
		return nil, fmt.Errorf("error getting text tracks: %w", err)
	}
	if result["error"] != nil {
		return []string{}, nil
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

func DumpVideo(config *model.Config, videoId string) (string, error) {
	result, err := executeRequest(config, fmt.Sprintf("https://api.vimeo.com/videos/%s", videoId))
	if err != nil {
		return "", fmt.Errorf("error getting video: %w", err)
	}
	return dumpJson(result)
}

func GetVideo(config *model.Config, videoId string) (*model.Video, error) {
	result, err := executeRequest(config, fmt.Sprintf("https://api.vimeo.com/videos/%s", videoId))
	if err != nil {
		return nil, fmt.Errorf("error getting video: %w", err)
	}
	description := ""
	if result["description"] != nil {
		description = result["description"].(string)
	}
	name := "unknown"
	if result["name"] != nil {
		name = result["name"].(string)
	}
	return &model.Video{
		ID:          videoId,
		Title:       name,
		Description: description,
	}, nil
}

func DumpUserVideos(config *model.Config, userId string) (string, error) {
	result, err := executeRequest(config, fmt.Sprintf("https://api.vimeo.com/users/%s/videos", userId))
	if err != nil {
		return "", fmt.Errorf("error getting user videos: %w", err)
	}
	return dumpJson(result)
}

func GetUserVideos(config *model.Config, userId string) ([]string, error) {
	links := []string{}
	url := fmt.Sprintf("https://api.vimeo.com/users/%s/videos", userId)

	for {
		result, err := executeRequest(config, url)
		if err != nil {
			return links, fmt.Errorf("error getting user videos: %w", err)
		}
		data, ok := result["data"].([]interface{})
		if !ok {
			return links, fmt.Errorf("unexpected response format: missing 'data'")
		}
		for _, video := range data {
			if videoMap, ok := video.(map[string]interface{}); ok {
				if videoId, ok := videoMap["link"].(string); ok {
					links = append(links, videoId)
				}
			}
		}
		log.Printf("Found %d videos.", len(links))
		paging, ok := result["paging"].(map[string]interface{})
		if !ok {
			break // Assume no more pages if format is wrong
		}
		nextRaw, ok := paging["next"]
		if !ok || nextRaw == nil {
			break // No more pages
		}
		next, ok := nextRaw.(string)
		if !ok || next == "" {
			break
		}
		url = fmt.Sprintf("https://api.vimeo.com%s", next)
	}
	return links, nil
}

func dumpJson(result map[string]interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}
	return string(jsonBytes), nil
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
