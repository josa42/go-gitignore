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

	if line := normalize(p.line); line != "" {
		if isPattern(line) {
			return patternToRegex(line).MatchString(path)
		}

		if isRoot(line) {
			return "/"+path == line
		}

		if isFilePath(line) {
			return filename(path) == line
		}
	}

	return false
}

func filename(line string) string {
	p := strings.Split(line, "/")
	return p[len(p)-1]
}

func isFilePath(line string) bool {
	return !isPattern(line)
}

// See: https://git-scm.com/docs/gitignore
//
// The pattern foo/ will match a directory foo and paths underneath it, but will
// not match a regular file or a symbolic link foo (this is consistent with the
// way how pathspec works in general in Git)
func isRoot(line string) bool {
	return strings.Contains(line, "/")
}

func isPattern(line string) bool {
	return strings.Contains(line, "*")
}

func normalize(line string) string {
	if !strings.HasPrefix(line, "/") && isRoot(line) {
		return "/" + line
	}
	return line
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
