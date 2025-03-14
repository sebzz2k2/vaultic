package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func WriteToFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func SearchLastMatch(filename, key string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	pattern := fmt.Sprintf(`^%s:(.*)`, key)
	re := regexp.MustCompile(pattern)

	var lastMatch string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); matches != nil {
			lastMatch = matches[1] // Extract content after "key:"
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastMatch, nil
}
