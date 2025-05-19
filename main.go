package main

import (
	"flag"
	"fmt"
	"log"
	vimeo_client "vimeo-transcriber-client"
)

func main() {
	videoID := flag.String("videoid", "vid", "The video id for which we are looking for a transcription")
	flag.Parse()

	if *videoID == "vid" {
		log.Fatal("Video ID is required. Please provide it using -videoid flag")
	}
	fmt.Println("video id: ", *videoID)

	loaded_config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}
	log.Println("Configuration loaded successfully")

	user, err := vimeo_client.GetUser(loaded_config)
	if err != nil {
		log.Fatal("Could not get user:", err)
	}
	log.Printf("User: %s", user.Username)

	textTracks, err := vimeo_client.GetTextTracks(loaded_config, *videoID)
	if err != nil {
		log.Fatal("Could not get text tracks:", err)
	}
	for _, track := range textTracks {
		log.Printf("Text track: %s", track)
	}
}
