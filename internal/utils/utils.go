package utils

import "strings"

func Match(pattern, path string) (bool, []string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make([]string, 0, len(patternParts))
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			params = append(params, pathParts[i])
		} else if part != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}