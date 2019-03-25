package util

import (
	"regexp"
)

func FindNamedGroups(regex *regexp.Regexp, input string) map[string]string {
	groups := make(map[string]string)

	match := regex.FindStringSubmatch(input)
	if match == nil {
		return groups
	}

	for i, gn := range regex.SubexpNames() {
		if i > 0 && gn != "" {
			groups[gn] = match[i]
		}
	}

	return groups
}
