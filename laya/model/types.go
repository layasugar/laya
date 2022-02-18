package model

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("(_|-)([a-zA-Z]+)")

func ToCamelCase(str string) string {
	camel := re.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)

	return camel
}

func strSpecialDeal(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, ";", "")
	return s
}
