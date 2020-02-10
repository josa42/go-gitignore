package gitignore

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Gitignore struct {
	patterns []Pattern
}

func NewGitignoreFromString(str string) Gitignore {

	fmt.Println(str)

	patterns := []Pattern{}

	for _, l := range strings.Split(str, "\n") {
		l = strings.TrimSpace(l)

		if l != "" && !strings.HasPrefix(l, "#") {
			patterns = append(patterns, Pattern{line: l})
		}
	}

	return Gitignore{patterns: patterns}
}

func NewGitignoreFromFile(path string) (Gitignore, error) {
	b, err := ioutil.ReadFile(path)
	return NewGitignoreFromString(string(b)), err
}

func (g Gitignore) Match(path string) bool {
	for _, p := range g.patterns {
		if p.Match(path) {
			return true
		}
	}
	return false
}