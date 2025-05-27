package main

import (
	"errors"
	"os"

	model "vimeo-transcriber-model"

	"github.com/joho/godotenv"
)

func loadConfig() (*model.Config, error) {
	error := godotenv.Load()
	if error != nil {
		return nil, error
	}
	vimeoPersonalAccessToken := os.Getenv("VIMEO_PERSONAL_ACCESS_TOKEN")
	if vimeoPersonalAccessToken == "" {
		return nil, errors.New("VIMEO_PERSONAL_ACCESS_TOKEN is not set. Please set it in the environment variables or use a .env file")
	}
	return &model.Config{
		VimeoPersonalAccessToken: vimeoPersonalAccessToken,
	}, nil
}
