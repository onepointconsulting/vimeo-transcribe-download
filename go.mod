module vimeo-transcriber

go 1.24.3

require (
	github.com/joho/godotenv v1.5.1
	vimeo-transcriber-client v0.0.0-00010101000000-000000000000
	vimeo-transcriber-model v0.0.0-00010101000000-000000000000
	vtt-parser v0.0.0-00010101000000-000000000000
)

replace vimeo-transcriber-model => ./model

replace vimeo-transcriber-client => ./vimeo_client

replace vtt-parser => ./vtt_parser
