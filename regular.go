package models

import (
	"regexp"
)

func ModelsFind(text string, model string) []string {
	re := regexp.MustCompile(model)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match
	}
	return  match
}

// ModelsFindAll 正则查找
func ModelsFindAll(html string, model string) [][]string {
	re := regexp.MustCompile(model)
	match := re.FindAllStringSubmatch(html,-1)
	if len(match) > 1 {
		return match
	}
	return  match
}
