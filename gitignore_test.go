package gitignore_test

import (
	"os"
	"testing"

	"github.com/josa42/go-gitignore"
	"github.com/stretchr/testify/assert"
)

func cd(p string) func() {
	pwd, _ := os.Getwd()
	os.Chdir(p)
	return func() {
		os.Chdir(pwd)
	}
}

func TestNewGitignoreFromString(t *testing.T) {
	assert.NotNil(t, gitignore.NewGitignoreFromString(""))
}

func TestNewGitignoreFromFile(t *testing.T) {
	defer cd("testdata/empty")()

	gitgnore, err := gitignore.NewGitignoreFromFile(".gitignore")

	assert.NotNil(t, gitgnore)
	assert.Nil(t, err)
}

func TestNewGitignoreFromFile_notFound(t *testing.T) {
	defer cd("testdata/does-not-exist")()

	gitgnore, err := gitignore.NewGitignoreFromFile(".gitignore")

	assert.NotNil(t, gitgnore)
	assert.NotNil(t, err)
}

func TestGitignoreMatch_empty(t *testing.T) {
	defer cd("testdata/empty")()

	i, _ := gitignore.NewGitignoreFromFile(".gitignore")

	assert.False(t, i.Match("foo.txt"))
	assert.False(t, i.Match("bar/foo.txt"))
}

func TestGitignoreMatch_basic(t *testing.T) {
	defer cd("testdata/basic")()

	i, _ := gitignore.NewGitignoreFromFile(".gitignore")

	assert.True(t, i.Match("ignore.txt"))
	assert.True(t, i.Match("bar/ignore.txt"))
	assert.False(t, i.Match("include.txt"))
	assert.False(t, i.Match("bar/include.txt"))
	assert.False(t, i.Match("other.txt"))
	assert.False(t, i.Match("bar/other.txt"))
}

func TestGitignoreMatch_wildcard(t *testing.T) {
	defer cd("testdata/wildcard")()

	i, _ := gitignore.NewGitignoreFromFile(".gitignore")

	assert.True(t, i.Match("data/ignore.txt"))
	assert.False(t, i.Match("data/keep.txt"))
}
