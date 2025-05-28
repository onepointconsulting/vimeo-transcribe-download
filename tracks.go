package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	vimeo_client "vimeo-transcriber-client"
	model "vimeo-transcriber-model"
	vtt_parser "vtt-parser"
)

func downloadTracks(loaded_config *model.Config, videoID string, targetFolder string) {
	textTracks, err := vimeo_client.GetTextTracks(loaded_config, videoID)
	if err != nil {
		log.Fatal("Could not get text tracks:", err)
	}
	video, err := vimeo_client.GetVideo(loaded_config, videoID)
	videoName := "unknown"
	if err == nil {
		videoName = strings.ReplaceAll(video.Title, " ", "_")
	} else {
		log.Printf("Could not get video details: %s", err)
		return
	}
	log.Printf("Video: %s", videoName)
	if len(textTracks) == 0 {
		var message = fmt.Sprintf("No text tracks found for video \"%s\" (https://vimeo.com/%s)", videoName, videoID)
		log.Print(message)
		appendToFile(targetFolder+"/missing_tracks.txt", fmt.Sprintf("%s\n", message))
	} else {
		loopTracks(textTracks, videoID, videoName, targetFolder)
	}
}

func loopTracks(textTracks []string, videoID, videoName, targetFolder string) {
	for i, track := range textTracks {
		log.Printf("Downloading text track: %s", track)
		fileName := fmt.Sprintf("%s_%s_%d.vtt", videoID, videoName, i)
		transcriptionFile := filepath.Join(targetFolder, fileName)
		fileName, err := vimeo_client.GetVttFile(track, transcriptionFile)
		if err != nil {
			log.Printf("Could not get vtt file: %s", err)
		} else {
			log.Printf("Downloaded text track: %s to %s", track, transcriptionFile)
			transcriptionText, err := vtt_parser.ExtractText(fileName)
			if err != nil {
				log.Printf("Could not extract text: %s", err)
			} else {
				fileName := fmt.Sprintf("%s_%s_%d.txt", videoID, videoName, i)
				transcriptionTextFile := filepath.Join(targetFolder, fileName)
				os.WriteFile(transcriptionTextFile, []byte(transcriptionText), 0644)
			}
		}

	}
}

func appendToFile(filePath string, text string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		log.Printf("Error writing to file: %s", err)
	}
}
