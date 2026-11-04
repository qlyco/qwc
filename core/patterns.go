package core

import (
	"fmt"
	"path"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func validatePattern(pattern string) error {
	for _, segment := range strings.Split(pattern, res.PathSeparator) {
		_, err := path.Match(segment, "")

		if err != nil {
			return fmt.Errorf(res.InvalidPatternMessage, pattern)
		}
	}

	return nil
}

func matchPattern(pattern string, name string) bool {
	return matchSegments(strings.Split(pattern, res.PathSeparator), strings.Split(name, res.PathSeparator))
}

func matchSegments(patterns []string, names []string) bool {
	if len(patterns) == 0 {
		return len(names) == 0
	}

	if patterns[0] == res.AnySegments {
		for i := 0; i <= len(names); i++ {
			if matchSegments(patterns[1:], names[i:]) {
				return true
			}
		}

		return false
	}

	if len(names) == 0 {
		return false
	}

	matched, _ := path.Match(patterns[0], names[0])

	return matched && matchSegments(patterns[1:], names[1:])
}
