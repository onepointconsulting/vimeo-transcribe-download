package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	vimeo_client "vimeo-transcriber-client"
)

const (
	COMMAND_TRANSCRIBE    = "transcribe"
	COMMAND_CHECK         = "check"
	ARG_VIDEO_ID          = "videoid"
	DEFAULT_VIDEO_ID      = "vid"
	ARG_TARGET_FOLDER     = "targetfolder"
	ARG_PRINT_USER        = "printuser"
	DEFAULT_TARGET_FOLDER = ""
	ARG_USER_ID           = "userid"
	DEFAULT_USER_ID       = "user"
	ARG_SIMPLE_VIDEO_LIST = "simplevideolist"
	ARG_TRANSCRIBE_USER   = "transcribeuser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'download' or 'check' subcommands")
		os.Exit(1)
	}

	loaded_config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not get configuration:", err)
	}
	log.Println("Configuration loaded successfully")

	switch os.Args[1] {
	case COMMAND_TRANSCRIBE:
		downloadCmd := flag.NewFlagSet(COMMAND_TRANSCRIBE, flag.ExitOnError)
		videoID := downloadCmd.String(ARG_VIDEO_ID, DEFAULT_VIDEO_ID, "The video id for which we are looking for a transcription")
		targetFolder := downloadCmd.String(ARG_TARGET_FOLDER, DEFAULT_TARGET_FOLDER, "The target folder for the transcription")

		downloadCmd.Parse(os.Args[2:])

		if *videoID == DEFAULT_VIDEO_ID {
			log.Fatal("Video ID is required. Please provide it using -videoid flag")
		}
		fmt.Println("video id: ", *videoID)

		downloadTracks(loaded_config, *videoID, *targetFolder)

	case COMMAND_CHECK:
		checkCmd := flag.NewFlagSet(COMMAND_CHECK, flag.ExitOnError)
		printuser := checkCmd.Bool(ARG_PRINT_USER, false, "Check if the user is to be printed")
		checkCmd.Parse(os.Args[2:])

		if *printuser {
			user, err := vimeo_client.GetUser(loaded_config)
			if err != nil {
				log.Fatal("Could not get user:", err)
			}
			log.Printf("User ID: %s", user.ID)
			log.Printf("User: %s", user.Username)

		} else {
			result, err := vimeo_client.DumpUser(loaded_config)
			if err != nil {
				log.Fatal("Could not get user:", err)
			}
			// Write the result to a file
			os.WriteFile("user.json", []byte(result), 0644)
		}
	case "dumpvideo":
		dumpvideoCmd := flag.NewFlagSet("dumpvideo", flag.ExitOnError)
		videoID := dumpvideoCmd.String(ARG_VIDEO_ID, DEFAULT_VIDEO_ID, "The video id for which we are looking for a transcription")
		dumpvideoCmd.Parse(os.Args[2:])

		if *videoID == DEFAULT_VIDEO_ID {
			log.Fatal("Video ID is required. Please provide it using -videoid flag")
		}

		result, err := vimeo_client.DumpVideo(loaded_config, *videoID)
		if err != nil {
			log.Fatal("Could not get video:", err)
		}
		// Write the result to a file
		os.WriteFile(fmt.Sprintf("video_%s.json", *videoID), []byte(result), 0644)

	case "uservideos":
		uservideosCmd := flag.NewFlagSet("uservideos", flag.ExitOnError)
		userID := uservideosCmd.String(ARG_USER_ID, DEFAULT_USER_ID, "The user id for which we are looking for videos")
		simpleVideoList := uservideosCmd.Bool(ARG_SIMPLE_VIDEO_LIST, false, "If true, the video list will be a simple list of video ids")
		transcribeUser := uservideosCmd.Bool(ARG_TRANSCRIBE_USER, false, "If true, tries to download the transcriptions of all videos of a psecified user")
		targetFolder := uservideosCmd.String(ARG_TARGET_FOLDER, DEFAULT_TARGET_FOLDER, "The target folder for the transcription")
		uservideosCmd.Parse(os.Args[2:])

		if *userID == DEFAULT_USER_ID {
			log.Fatal("User ID is required. Please provide it using -userid flag")
		}

		if *simpleVideoList || *transcribeUser {
			videos, err := vimeo_client.GetUserVideos(loaded_config, *userID)
			if err != nil {
				log.Fatal("Could not get user videos:", err)
			}
			if *simpleVideoList {
				// Write the result to a file
				os.WriteFile(fmt.Sprintf("user_videos_%s.txt", *userID), []byte(strings.Join(videos, "\n")), 0644)
				log.Printf("Wrote %d user videos to user_videos_%s.txt", len(videos), *userID)
			} else {
				// Write the result to a file
				for i, video := range videos {
					parsedUrl, err := url.Parse(video)
					if err != nil {
						log.Printf("Could not parse video url: %s", err)
					} else {
						videoID := parsedUrl.Path[1:]
						downloadTracks(loaded_config, videoID, *targetFolder)
					}
					log.Printf("Downloaded video %d of %d", i+1, len(videos))
				}
			}
		} else {
			result, err := vimeo_client.DumpUserVideos(loaded_config, *userID)
			if err != nil {
				log.Fatal("Could not get user videos:", err)
			}
			// Write the result to a file
			os.WriteFile(fmt.Sprintf("user_videos_%s.json", *userID), []byte(result), 0644)
		}

	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		fmt.Println("expected 'download' or 'check' subcommands")
		os.Exit(1)
	}
}
