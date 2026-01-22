package exporter

import "strings"

func containsAny(s string, parts ...string) bool {
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
