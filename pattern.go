package gitignore

import (
	"fmt"
	"regexp"
	"strings"
)

type Pattern struct {
	line string
}

var (
	exprStripPrefixSlash    = regexp.MustCompile("^/")
	exprStripSuffixSlash    = regexp.MustCompile("/$")
	exprSplitPattern        = regexp.MustCompile(`\*\*(\/\*)?`)
	exprQuotedCharacterList = regexp.MustCompile(`\\\[([^\]]*)\\\]`)
	exprQuotedSpace         = regexp.MustCompile(`\\ `)
)

func (p Pattern) Match(path string) bool {

	if line := normalize(p.line); line != "" {

		// fmt.Printf("line = '%s'\npath = '%s'\n", p.line, path)
		return patternToRegex(line).MatchString(path)
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

func isDirPath(line string) bool {
	return strings.HasSuffix(line, "/")
}

// See: https://git-scm.com/docs/gitignore
//
// The pattern foo/ will match a directory foo and paths underneath it, but will
// not match a regular file or a symbolic link foo (this is consistent with the
// way how pathspec works in general in Git)
func isRootPath(line string) bool {
	return strings.Contains(stripSuffixSlash(line), "/")
}

func isPattern(line string) bool {
	return true
}

func normalize(line string) string {
	line = strings.ReplaceAll(line, "\\#", "#")

	quotedSpace := strings.HasSuffix(line, "\\ ")
	// FIXME
	line = strings.TrimSpace(line)
	if quotedSpace {
		line += " "
	}
	if !strings.HasPrefix(line, "/") && isRootPath(line) {
		return "/" + line
	}
	return line
}

func patternToRegex(pattern string) *regexp.Regexp {

	prefix := ""
	if strings.HasPrefix(pattern, "/") {
		prefix = "^"
		pattern = pattern[1:]
	}

	// root := isRootPath(pattern)
	dir := isDirPath(pattern)
	childrenOnly := strings.HasSuffix(pattern, "/**") || strings.HasSuffix(pattern, "/*")

	if strings.Contains(pattern, "**") && !strings.HasPrefix(pattern, "**") {
		prefix = "^"
	}

	pat := []string{}
	for _, p := range exprSplitPattern.Split(pattern, -1) {
		innerPat := []string{}

		for _, pi := range strings.Split(p, "*") {
			innerPat = append(innerPat, quotePattern(pi))
		}

		pat = append(pat, strings.Join(innerPat, `[^/]*`))
	}

	exprStr := prefix + strings.Join(pat, `(.*|)`)

	if childrenOnly {
		exprStr += "(/.*)"
	} else if dir {
		exprStr += "(/.*)"
	} else {
		exprStr += "($|/)"
	}

	// fmt.Println("expr =", exprStr)
	exp, _ := regexp.Compile(exprStr)

	return exp
}

func quotePattern(str string) string {
	str = stripSurroundingSlashes(str)
	str = regexp.QuoteMeta(str)

	str = strings.ReplaceAll(str, "\\\\ ", "\\s")

	// do not quote character lists like [abc]
	if m := exprQuotedCharacterList.FindAllStringSubmatch(str, -1); m != nil {
		for _, find := range m {
			str = strings.Replace(str, find[0], fmt.Sprintf(`[%s]+`, find[1]), 1)
		}
	}

	return str
}

func stripSuffixSlash(str string) string {
	return exprStripSuffixSlash.ReplaceAllString(str, "")
}

func stripPrefixSlash(str string) string {
	return exprStripPrefixSlash.ReplaceAllString(str, "")
}

func stripSurroundingSlashes(str string) string {
	return stripPrefixSlash(stripSuffixSlash(str))
}

func ensureSuffixSlash(str string) string {
	return exprStripSuffixSlash.ReplaceAllString(str, "") + "/"
}
