package main

import (
	"regexp"
)

func GetLinks(body string) []string {
	var links []string
	set := make(map[string]struct{})

	re := regexp.MustCompile(`href="(https?://[^"]*)"`)
	matches := re.FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		if _, exists := set[match[1]]; !exists {
			links = append(links, match[1])
			set[match[1]] = struct{}{}
		}
	}

	return links	
}