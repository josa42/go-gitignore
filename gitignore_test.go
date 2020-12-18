package gitignore_test

import (
	"encoding/json"
	"io/ioutil"
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

type Case struct {
	Name     string
	Ignore   string `json:"ignore"`
	FilePath string `json:"file_path"`
	Ignored  bool   `json:"ignored"`
}

func TestGitignore(t *testing.T) {

	cases := []Case{}

	content, _ := ioutil.ReadFile("testdata/cases.json")
	json.Unmarshal(content, &cases)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ignore := gitignore.NewGitignoreFromString(c.Ignore)

			assert.Equal(t, c.Ignored, ignore.Match(c.FilePath))
		})
	}

}
