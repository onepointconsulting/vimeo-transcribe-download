# Vimeo Transcriber

Simple command line tool which can be used to download the transcriptions on paid accounts from Vimeo.

## Building

Use:

`go build .`

## Configuration

You should define this environment variable:

```
VIMEO_PERSONAL_ACCESS_TOKEN=<token>
```

You need to get this token from Vimeo.

## Running

### Transcription of single video

Examples:

```
vimeo-transcriber.exe transcribe -videoid 1085658122 -targetfolder transcriptions
vimeo-transcriber.exe transcribe -videoid 1085780999 -targetfolder transcriptions
```

This will download the vtt files to the `transcriptions` folder in case there is one.

```
2025/05/19 19:19:37 Text track: https://captions.cloud.vimeo.com/captions/231961801.vtt?expires=1747685972&sig=72561a1b6f7a5f519312a4e79872fc6b102d0ef5&download=auto_generated_captions.vtt&hls=1
```

If there is no text track, you will see a message like this one:

```
2025/05/27 09:42:19 No text tracks found for video 1085658122
```

### User Details

Examples:

```
vimeo-transcriber.exe check -printuser
vimeo-transcriber.exe check
```

### Dumping video information

Examples:

```
vimeo-transcriber.exe dumpvideo -videoid 468121987
```

### Fetching all of the videos of a user

Examples:

```
vimeo-transcriber.exe uservideos -userid user51871371 -simplevideolist
```
