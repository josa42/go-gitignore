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
	exprSpaceSuffix         = regexp.MustCompile(`\\ +$`)
)

func (p Pattern) Match(path string) bool {

	if line := normalize(p.line); line != "" {
		return patternToRegex(line).MatchString(path)
	}

	return false
}

// See: https://git-scm.com/docs/gitignore
//
// The pattern foo/ will match a directory foo and paths underneath it, but will
// not match a regular file or a symbolic link foo (this is consistent with the
// way how pathspec works in general in Git)
func isRootPath(line string) bool {
	return strings.HasPrefix(line, "/") || strings.Contains(stripSuffixSlash(line), "/")
}

func normalize(line string) string {
	line = strings.ReplaceAll(line, "\\#", "#")

	keepSpaceSuffix := exprSpaceSuffix.MatchString(line)
	line = strings.TrimSpace(line)
	if keepSpaceSuffix {
		line += " "
	}

	if isRootPath(line) && !strings.HasPrefix(line, "/") {
		line = "/" + line
	}

	return line
}

func patternToRegex(pattern string) *regexp.Regexp {

	prefix := ""
	if strings.HasPrefix(pattern, "/") {
		prefix = "^"
		pattern = pattern[1:]
	}

	childrenOnly := strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "/**") || strings.HasSuffix(pattern, "/*")

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
	} else {
		exprStr += "($|/)"
	}

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
