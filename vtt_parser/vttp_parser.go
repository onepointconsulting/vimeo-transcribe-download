package vtt_parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ExtractText(vttFile string) (string, error) {
	file, err := os.Open(vttFile)
	if err != nil {
		log.Printf("Error reading vtt file: %s", err)
		return "", fmt.Errorf("error reading vtt file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	texts := []string{}
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip header lines (typically first 3 lines in VTT)
		if lineNumber <= 3 || strings.Contains(line, " --> ") || len(line) == 0 {
			continue
		}

		_, err := strconv.Atoi(line)
		if err == nil {
			// It is a number, skip it
			continue
		}
		texts = append(texts, line)
	}

	return strings.Join(texts, " "), nil
}
