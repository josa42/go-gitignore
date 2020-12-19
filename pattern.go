package gitignore

import (
	"regexp"
	"strings"
)

type Pattern struct {
	line string
}

var (
	exprStripSuffixSlash = regexp.MustCompile("/$")
)

func (p Pattern) Match(path string) bool {

	if line := normalize(p.line); line != "" {
		if isPattern(line) {
			return patternToRegex(line).MatchString(path)
		}

		if isRoot(line) {
			return "/"+path == line || strings.HasPrefix("/"+path, ensureSuffixSlash(line))
		}

		if isFilePath(line) {
			return filename(path) == line
		}

		if isDirPath(line) {
			return dirname(path) == dirname(line)
		}
	}

	return false
}

func filename(line string) string {
	p := strings.Split(line, "/")
	return p[len(p)-1]
}

func dirname(line string) string {
	p := strings.Split(line, "/")
	if len(p) >= 2 {
		return p[len(p)-2]
	}
	return ""
}

func isFilePath(line string) bool {
	return !isPattern(line) && !strings.HasSuffix(line, "/")
}

func isDirPath(line string) bool {
	return !isPattern(line) && strings.HasSuffix(line, "/")
}

// See: https://git-scm.com/docs/gitignore
//
// The pattern foo/ will match a directory foo and paths underneath it, but will
// not match a regular file or a symbolic link foo (this is consistent with the
// way how pathspec works in general in Git)
func isRoot(line string) bool {
	return strings.Contains(stripSuffixSlash(line), "/")
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

	exp, _ := regexp.Compile(prefix + strings.Join(pat, `.*`))

	return exp
}

func stripSuffixSlash(str string) string {
	return exprStripSuffixSlash.ReplaceAllString(str, "")
}

func ensureSuffixSlash(str string) string {
	return exprStripSuffixSlash.ReplaceAllString(str, "") + "/"
}
