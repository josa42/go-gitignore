package gitignore

import (
	"io/ioutil"
	"strings"
)

type Gitignore struct {
	patterns []Pattern
	excludes []Pattern
}

func NewGitignoreFromString(str string) Gitignore {

	patterns := []Pattern{}
	excludes := []Pattern{}

	for _, l := range strings.Split(str, "\n") {
		l = strings.TrimSpace(l)

		if l != "" && !strings.HasPrefix(l, "#") {
			if strings.HasPrefix(l, "!") {
				excludes = append(excludes, Pattern{line: l[1:]})
			} else {
				patterns = append(patterns, Pattern{line: l})
			}
		}
	}

	return Gitignore{patterns: patterns, excludes: excludes}
}

func NewGitignoreFromFile(path string) (Gitignore, error) {
	b, err := ioutil.ReadFile(path)
	return NewGitignoreFromString(string(b)), err
}

func (g Gitignore) Match(path string) bool {
	for _, p := range g.excludes {
		if p.Match(path) {
			return false
		}
	}
	for _, p := range g.patterns {
		if p.Match(path) {
			return true
		}
	}
	return false
}
