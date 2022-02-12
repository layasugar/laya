package template

import "strings"

func parseName(name string) (string, string) {
	ns := strings.Split(name, "/")
	if len(ns) > 0 {
		return ns[len(ns)-1], name
	}
	return name, name
}
