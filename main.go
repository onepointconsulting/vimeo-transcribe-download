package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	vimeo_client "vimeo-transcriber-client"
)

const (
	ARG_VIDEO_ID          = "videoid"
	DEFAULT_VIDEO_ID      = "vid"
	ARG_TARGET_FOLDER     = "targetfolder"
	ARG_CHECK_USER        = "checkuser"
	DEFAULT_TARGET_FOLDER = ""
)

func main() {
	videoID := flag.String(ARG_VIDEO_ID, DEFAULT_VIDEO_ID, "The video id for which we are looking for a transcription")
	targetFolder := flag.String(ARG_TARGET_FOLDER, DEFAULT_TARGET_FOLDER, "The target folder for the transcription")
	checkUser := flag.Bool(ARG_CHECK_USER, false, "Check if the user is authenticated")
	flag.Parse()

	if *videoID == DEFAULT_VIDEO_ID {
		log.Fatal("Video ID is required. Please provide it using -videoid flag")
	}
	fmt.Println("video id: ", *videoID)

	loaded_config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}
	log.Println("Configuration loaded successfully")

	if *checkUser {
		user, err := vimeo_client.GetUser(loaded_config)
		if err != nil {
			log.Fatal("Could not get user:", err)
		}
		log.Printf("User: %s", user.Username)
	}

	textTracks, err := vimeo_client.GetTextTracks(loaded_config, *videoID)
	if err != nil {
		log.Fatal("Could not get text tracks:", err)
	}
	if len(textTracks) == 0 {
		log.Fatalf("No text tracks found for video %s", *videoID)
	}
	for _, track := range textTracks {
		log.Printf("Downloading text track: %s", track)
		fileName := fmt.Sprintf("%s.vtt", *videoID)
		transcriptionFile := filepath.Join(*targetFolder, fileName)
		vimeo_client.GetVttFile(track, transcriptionFile)
		log.Printf("Downloaded text track: %s to %s", track, transcriptionFile)
	}
}
