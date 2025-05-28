package main

import (
	"errors"
	"os"
	"strconv"

	model "vimeo-transcriber-model"

	"github.com/joho/godotenv"
)

func loadConfig() (*model.Config, error) {
	error := godotenv.Load()
	if error != nil {
		return nil, error
	}
	vimeoPersonalAccessToken := os.Getenv("VIMEO_PERSONAL_ACCESS_TOKEN")
	retryDelayString := os.Getenv("RETRY_DELAY")
	if vimeoPersonalAccessToken == "" {
		return nil, errors.New("VIMEO_PERSONAL_ACCESS_TOKEN is not set. Please set it in the environment variables or use a .env file")
	}
	if retryDelayString == "" {
		retryDelayString = "30"
	}
	retryDelay, err := strconv.Atoi(retryDelayString)
	if err != nil {
		return nil, err
	}
	return &model.Config{
		VimeoPersonalAccessToken: vimeoPersonalAccessToken,
		RetryDelay:               retryDelay,
	}, nil
}
