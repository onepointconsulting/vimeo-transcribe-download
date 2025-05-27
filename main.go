package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	vimeo_client "vimeo-transcriber-client"
)

const (
	ARG_VIDEO_ID          = "videoid"
	DEFAULT_VIDEO_ID      = "vid"
	ARG_TARGET_FOLDER     = "targetfolder"
	ARG_PRINT_USER        = "printuser"
	DEFAULT_TARGET_FOLDER = ""
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'download' or 'check' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "download":
		downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
		videoID := downloadCmd.String(ARG_VIDEO_ID, DEFAULT_VIDEO_ID, "The video id for which we are looking for a transcription")
		targetFolder := downloadCmd.String(ARG_TARGET_FOLDER, DEFAULT_TARGET_FOLDER, "The target folder for the transcription")

		downloadCmd.Parse(os.Args[2:])

		if *videoID == DEFAULT_VIDEO_ID {
			log.Fatal("Video ID is required. Please provide it using -videoid flag")
		}
		fmt.Println("video id: ", *videoID)

		loaded_config, err := loadConfig()
		if err != nil {
			log.Fatal("Could not get configuration:", err)
		}
		log.Println("Configuration loaded successfully")

		textTracks, err := vimeo_client.GetTextTracks(loaded_config, *videoID)
		if err != nil {
			log.Fatal("Could not get text tracks:", err)
		}
		if len(textTracks) == 0 {
			log.Fatalf("No text tracks found for video %s", *videoID)
		}
		loopTracks(textTracks, *videoID, *targetFolder)

	case "check":
		checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
		printuser := checkCmd.Bool(ARG_PRINT_USER, true, "Check if the user is to be printed")
		checkCmd.Parse(os.Args[2:])

		loaded_config, err := loadConfig()
		if err != nil {
			log.Fatal("Could not get configuration:", err)
		}
		log.Println("Configuration loaded successfully")

		if *printuser {
			user, err := vimeo_client.GetUser(loaded_config)
			if err != nil {
				log.Fatal("Could not get user:", err)
			}
			log.Printf("User ID: %s", user.ID)
			log.Printf("User: %s", user.Username)
		}

	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		fmt.Println("expected 'download' or 'check' subcommands")
		os.Exit(1)
	}
}

func loopTracks(textTracks []string, videoID string, targetFolder string) {
	for _, track := range textTracks {
		log.Printf("Downloading text track: %s", track)
		fileName := fmt.Sprintf("%s.vtt", videoID)
		transcriptionFile := filepath.Join(targetFolder, fileName)
		vimeo_client.GetVttFile(track, transcriptionFile)
		log.Printf("Downloaded text track: %s to %s", track, transcriptionFile)
	}
}
