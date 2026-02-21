package main

import (
	"regexp"
)

func GetLinks(body string) []string {
	var links []string

	re := regexp.MustCompile(`href="(https?://[^"]*)"`)
	matches := re.FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		links = append(links, match[1])
	}

	return links	
}