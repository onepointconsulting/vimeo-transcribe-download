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

Example:

```
vimeo-transcriber.exe -videoid 1085658122 -targetfolder transcriptions
vimeo-transcriber.exe -videoid 1085780999 -targetfolder transcriptions
```

This will download the vtt files to the `transcriptions` folder.

```
2025/05/19 19:19:37 Text track: https://captions.cloud.vimeo.com/captions/231961801.vtt?expires=1747685972&sig=72561a1b6f7a5f519312a4e79872fc6b102d0ef5&download=auto_generated_captions.vtt&hls=1
```