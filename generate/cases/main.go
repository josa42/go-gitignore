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
	Skip     bool   `json:"skip"`
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	FilePath string `json:"file_path"`
	Ignored  bool   `json:"ignored"`
}

type Repo struct {
	Skip    bool     `json:"skip"`
	Name    string   `json:"name"`
	Pattern []string `json:"pattern"`
	Files   []string `json:"files"`
}

// Test from:
// - https://github.com/zabawaba99/go-gitignore/blob/master/ignore_test.go
// - https://github.com/nathankleyn/gitignore.rs/tree/master/tests/resources/fake_repo
// - https://github.com/mherrmann/gitignore_parser/blob/master/tests.py
// - https://github.com/hawknewton/gitignore-parser/blob/master/spec/gitignore/parser/rule_spec.rb

func main() {
	repos := []Repo{}
	content, _ := ioutil.ReadFile("generate/repos.json")
	json.Unmarshal(content, &repos)

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

	pattern := strings.Join(repo.Pattern, "\n")
	ioutil.WriteFile(".gitignore", []byte(pattern), 07755)

	for _, f := range repo.Files {
		cases = append(cases, Case{
			Skip:     repo.Skip,
			Name:     fmt.Sprintf("%s: %s", repo.Name, f),
			Pattern:  pattern,
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
	out, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			return false
		}

		log.Fatal(string(out))
	}

	return true
}
