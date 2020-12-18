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
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	FilePath string `json:"file_path"`
	Ignored  bool   `json:"ignored"`
}

type Repo struct {
	Name    string
	Pattern []string
	Files   []string
}

var repos = []Repo{
	{
		Name:    "Simple",
		Pattern: []string{"/data"},
		Files:   []string{"data", "data/file", "data.json"},
	},
	{
		Name:    "basic",
		Pattern: []string{"ignore.txt"},
		Files:   []string{"ignore.txt", "bar/ignore.txt", "include.txt", "bar/include.txt", "other.txt", "bar/other.txt"},
	},
	{
		Name:    "empty",
		Pattern: []string{},
		Files:   []string{"foo.txt", "bar/foo.txt"},
	},
	{
		Name:    "Wildcards",
		Pattern: []string{"data/*", "!data/keep.txt", "/rootdata/*", "!keep.all.txt", "tmp/"},
		Files: []string{
			"data/ignore.txt",
			"data/keep.txt",
			"rootdata/data/ignore.txt",
			"data/data/ignore.txt",
			"other/data/ignore.txt",
			"other/data/keep.txt",
			"tmp",
			"tmp/file",
			"sub/tmp",
			"sub/tmp/file",
		},
	},
	// https://github.com/zabawaba99/go-gitignore/blob/master/ignore_test.go
	{
		Name:    "star matching",
		Pattern: []string{"/*.txt"},
		Files:   []string{"foo.txt", "somedir/foo.txt"},
	},
	{
		Name:    "double star prefix",
		Pattern: []string{"**/foo.txt"},
		Files:   []string{"hello/foo.txt", "some/dirs/foo.txt"},
	},
	{
		Name:    "double star suffix",
		Pattern: []string{"/hello/**"},
		Files:   []string{"hello/foo.txt", "some/dirs/foo.txt"},
	},
	{

		Name:    "double star in path",
		Pattern: []string{"/hello/**/world.txt"},
		Files:   []string{"hello/world.txt", "hello/stuff/world.txt", "some/dirs/foo.txt"},
	},
	{
		Name:    "negate doubl start patterns",
		Pattern: []string{"!**/foo.txt"},
		Files:   []string{"hello/foo.txt", "hello/foo.txt", "hello/world.txt"},
	},
	// https://github.com/nathankleyn/gitignore.rs/tree/master/tests/resources/fake_repo
	{
		Name: "gitignore.rs",
		Pattern: []string{
			"*.no",
			"not_me_either/",
			"/or_even_me",
		},
		Files: []string{
			"a_dir/a_nested_dir/deeper_still/hola.greeting",
			"a_dir/a_nested_dir/deeper_still/hello.greeting",
			"a_dir/a_nested_dir/.gitignore",
			"also_include_me",
			"include_me",
			"not_me.no",
			"not_me_either/i_shouldnt_be_included",
			"or_even_me",
			"or_me.no",
		},
	},

	// https://github.com/mherrmann/gitignore_parser/blob/master/tests.py
	{
		Name: "mherrmann/gitignore_parser: simple",
		Pattern: []string{
			"__pycache__/\n",
			"*.py[cod]",
		},
		Files: []string{
			"main.py",
			"main.pyc",
			"dir/main.pyc",
			"__pycache__",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: wildcard",
		Pattern: []string{
			"hello.*",
		},
		Files: []string{
			"hello.txt",
			"hello.foobar/",
			"dir/hello.txt",
			"hello.",
			"hello",
			"helloX",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: anchored_wildcard",
		Pattern: []string{
			"/hello.*",
		},
		Files: []string{
			"hello.txt",
			"hello.c",
			"a/hello.java",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: trailingspaces",
		Pattern: []string{
			"ignoretrailingspace ",
			"notignoredspace\\ ",
			"partiallyignoredspace\\  ",
			"partiallyignoredspace2 \\  ",
			"notignoredmultiplespace\\ \\ \\ ",
		},
		Files: []string{
			"ignoretrailingspace",
			"ignoretrailingspace ",
			"partiallyignoredspace ",
			"partiallyignoredspace  ",
			"partiallyignoredspace",
			"partiallyignoredspace2  ",
			"partiallyignoredspace2   ",
			"partiallyignoredspace2 ",
			"partiallyignoredspace2",
			"notignoredspace ",
			"notignoredspace",
			"notignoredmultiplespace   ",
			"notignoredmultiplespace",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: comment",
		Pattern: []string{
			"somematch\n",
			"#realcomment\n",
			"othermatch\n",
			"\\#imnocomment",
		},
		Files: []string{
			"somematch",
			"#realcomment",
			"othermatch",
			"#imnocomment",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: ignore_directory",
		Pattern: []string{
			".venv/",
		},
		Files: []string{
			".venv",
			".venv/folder",
			".venv/file.txt",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: directory_asterisk",
		Pattern: []string{
			".venv/*",
		},
		Files: []string{
			".venv",
			".venv/folder",
			".venv/file.txt",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: negation",
		Pattern: []string{
			"*.ignore",
			"!keep.ignore",
		},
		Files: []string{
			"trash.ignore",
			"keep.ignore",
			"waste.ignore",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: double_asterisks",
		Pattern: []string{
			"foo/**/Bar",
		},
		Files: []string{
			"foo/hello/Bar",
			"foo/world/Bar",
			"foo/Bar",
		},
	},
	{
		Name: "mherrmann/gitignore_parser: single_asterisk",
		Pattern: []string{
			"*",
		},
		Files: []string{
			"file.txt",
			"directory",
			"directory-trailing/",
		},
	},

	// https://github.com/hawknewton/gitignore-parser/blob/master/spec/gitignore/parser/rule_spec.rb
	{
		Name:    "hawknewton/gitignore-parser: given a patern without slashes",
		Pattern: []string{"foo*bar"},
		Files: []string{
			"foobar",
			"foo123bar",
			"bar",
			"foobar/",
			"foobar123/",
			"foobar",
			"test",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with unescaped trailing spaces",
		Pattern: []string{"foo*  "},
		Files: []string{
			"foo",
			"foo/",
			"bar/",
			"foobar",
			"test",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with escaped trailing spaces",
		Pattern: []string{"foo*\\ \\ "},
		Files: []string{
			"foo123  ",
			"foo123 ",
			"foo123  /",
			"foo123 /",
			"foobar",
			"test",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading #",
		Pattern: []string{"#foo"},
		Files: []string{
			"#foo",
			"foo",
			"foobar",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading #",
		Pattern: []string{"#foo"},
		Files: []string{
			"#foo",
			"foo",
			"#foo/",
			"foo/",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern ending with a slash",
		Pattern: []string{"foo*/"},
		Files: []string{
			"foo",
			"foo/",
			"foo123/",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern starting with **",
		Pattern: []string{"**/foo"},
		Files: []string{
			"foo",
			"bar",
			"foo/",
			"bar/",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading slash",
		Pattern: []string{"/foo"},
		Files: []string{
			"foo",
			"bar",
			"foo/",
			"bar/",
			"bar",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with multiple slashes",
		Pattern: []string{"foo/bar/baz"},
		Files: []string{
			"foo",
			"foo/",
			"baz",
			"baz/",
			"test",
			"bar/baz",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading slash and multiple slashes",
		Pattern: []string{"/foo/bar/baz"},
		Files: []string{
			"foo",
			"foo/",
			"baz",
			"baz/",
			"test",
			"bar/baz",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with multiple slashes and two asterisks",
		Pattern: []string{"foo/**/baz"},
		Files: []string{
			"foo",
			"foo/",
			"baz",
			"baz/",
			"test",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading slash and two asterisks",
		Pattern: []string{"/**/baz"},
		Files: []string{
			"foo",
			"foo/",
			"baz",
			"baz/",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a a trailing \"/**\"",
		Pattern: []string{"baz/**"},
		Files: []string{
			"foo",
			"baz",
			"baz/",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern of \" * *\"",
		Pattern: []string{"**"},
		Files: []string{
			"foo",
			"baz/",
			"test",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern without a leading !",
		Pattern: []string{"foo"},
		Files: []string{
			"foo",
		},
	},
	{
		Name:    "hawknewton/gitignore-parser: given a pattern with a leading !",
		Pattern: []string{"!foo"},
		Files: []string{
			"foo",
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

	pattern := strings.Join(repo.Pattern, "\n")
	ioutil.WriteFile(".gitignore", []byte(pattern), 07755)

	for _, f := range repo.Files {
		cases = append(cases, Case{
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
