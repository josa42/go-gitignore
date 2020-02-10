package gitignore

import (
	"fmt"
	"regexp"
	"strings"
)

type Pattern struct {
	line string
}

func (p Pattern) Match(path string) bool {

	if p.line != "" {
		if isFilename(p.line) {
			return filename(path) == p.line
		}

		if isRoot(p.line) && !isPattern(p.line) {
			return "/"+path == p.line
		}

		if isPattern(p.line) {
			return patternToRegex(p.line).MatchString(path)
		}

		// TODO handle wildcards
	}

	return false
}

func filename(line string) string {
	p := strings.Split(line, "/")
	return p[len(p)-1]
}

func isFilename(line string) bool {
	return !strings.Contains(line, "/") && !isPattern(line)
}

func isRoot(line string) bool {
	return strings.HasPrefix(line, "/")
}

func isPattern(line string) bool {
	return strings.Contains(line, "*")
}

var splitExp = regexp.MustCompile(`\*\*(\/\*)?`)

func patternToRegex(pattern string) *regexp.Regexp {

	prefix := ""
	if strings.HasPrefix(pattern, "/") {
		prefix = "^"
		pattern = pattern[1:]
	}

	if strings.Contains(pattern, "**") && !strings.HasPrefix(pattern, "**") {
		prefix = "^"
	}

	pat := []string{}
	for _, p := range splitExp.Split(pattern, -1) { // strings.Split(patter, "**") {
		innerPat := []string{}
		for _, pi := range strings.Split(p, "*") {
			innerPat = append(innerPat, regexp.QuoteMeta(pi))
		}
		pat = append(pat, strings.Join(innerPat, `[^/]*`))
	}

	fmt.Printf("p: %s\n", prefix+strings.Join(pat, `.*`))

	exp, _ := regexp.Compile(prefix + strings.Join(pat, `.*`))

	return exp
}
