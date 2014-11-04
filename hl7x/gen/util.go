package main

import (
	"strings"
	"unicode"

	"bitbucket.org/pkg/inflect"
)

func NormalizeTypeName(rs *inflect.Ruleset, v string) string {
	if len(v) == 0 {
		return ""
	}

	if unicode.IsLower(rune(v[0])) {
		return rs.Capitalize(v)
	}
	return rs.Capitalize(rs.Camelize(strings.ToLower(v)))
}

func NormalizeFileName(rs *inflect.Ruleset, v string) string {
	if len(v) == 0 {
		return ""
	}

	if unicode.IsLower(rune(v[0])) {
		return rs.Underscore(v)
	}
	return strings.ToLower(v)
}
