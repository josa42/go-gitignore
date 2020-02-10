package gitignore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPattern_Match(t *testing.T) {
	tests := []struct {
		line    string
		pattern string
		expect  bool
	}{
		{"", "", false},
		{"foo.txt", "", false},
		{"foo.txt", "foo.txt", true},
		{"foo/foo.txt", "foo.txt", true},
		{"foo.txt", "/foo.txt", true},
		{"foo/foo.txt", "/foo.txt", false},
		{"bin", "/bin", true},
		{"vendor/bin", "/bin", false},
		{"foo.txt", "*.txt", true},
		{"foo/foo.txt", "*.txt", true},
		{"foo.txt", "/*.txt", true},
		{"foo/foo.txt", "/*.txt", false},
		{".git", ".git*", true},
		{".gitignore", ".git*", true},
		{"foo.txt", "**/*.txt", true},
		{"foo.txt", "/**/*.txt", true},
		{"foo.txt", "/bar/**/*.txt", false},
		{"bar/foo/test.txt", "/foo/**/*.txt", false},
		{"foo/bar/test.txt", "/foo/**/*.txt", true},
		{"foo/test.txt", "/foo/**/*.txt", true},
		{"bar/foo/test.txt", "foo/**/*.txt", false},
		{"foo/bar/test.txt", "foo/**/*.txt", true},
		{"foo/test.txt", "foo/**/*.txt", true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expect, Pattern{line: tt.pattern}.Match(tt.line), fmt.Sprintf(`"%s" => "%s"`, tt.line, tt.pattern))
		})
	}
}
