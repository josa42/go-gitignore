package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	Name     string
	Ignore   string `json:"ignore"`
	FilePath string `json:"file_path"`
	Ignored  bool   `json:"ignored"`
}

type Repo struct {
	Name   string
	Ignore []string
	Files  []string
}

var repos = []Repo{
	{
		Name:   "Simple",
		Ignore: []string{"/data"},
		Files:  []string{"data", "data/file", "data.json"},
	},
	{
		Name:   "basic",
		Ignore: []string{"ignore.txt"},
		Files:  []string{"ignore.txt", "bar/ignore.txt", "include.txt", "bar/include.txt", "other.txt", "bar/other.txt"},
	},
	{
		Name:   "empty",
		Ignore: []string{},
		Files:  []string{"foo.txt", "bar/foo.txt"},
	},
	{
		Name:   "Wildcards",
		Ignore: []string{"data/*", "!data/keep.txt", "/rootdata/*", "!keep.all.txt", "tmp/"},
		Files: []string{
			"data/ignore.txt",
			"data/keep.txt",
			"rootdata/data/ignore.txt",
			"data/data/ignore.txt",
			"other/data/ignore.txt",
			"other/data/keep.txt",
			"tmp",
			"tmp/file",
		},
	},
}

func main() {
	cases := []Case{}

	for _, repo := range repos {
		cases = append(cases, generate(repo)...)
	}

	s, _ := json.MarshalIndent(cases, "", "  ")
	err := ioutil.WriteFile("testdata/cases.json", []byte(s), 0755)
	if err != nil {
		log.Fatal(err)
	}

}

func generate(repo Repo) []Case {
	cases := []Case{}

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	dir, _ := ioutil.TempDir("", "generate-test-tmp-*")
	defer os.RemoveAll(dir)
	os.Chdir(dir)

	run("git", "init")

	ignore := strings.Join(repo.Ignore, "\n")
	ioutil.WriteFile(".gitignore", []byte(ignore), 07755)

	for _, f := range repo.Files {
		cases = append(cases, Case{
			Name:     fmt.Sprintf("%s: %s", repo.Name, f),
			Ignore:   ignore,
			FilePath: f,
			Ignored:  ignored(f),
		})

	}

	return cases
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(out))
		log.Fatal(err)
	}
}

func ignored(filePath string) bool {
	cmd := exec.Command("git", "check-ignore", filePath)
	err := cmd.Run()
	if err != nil {
		if err.Error() == "exit status 1" {
			return false
		}
		log.Fatal(err)
	}

	return true
}
