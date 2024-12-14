package nerfs

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-set/v3"
)

const (
	DomainsFile = "domains.txt"
	WordsFile   = "wordlist.json"
)

type Artifact struct {
	domains *set.Set[string]
	words   []*regexp.Regexp
}

func Load(directory string) (*Artifact, error) {
	domainsFile := filepath.Join(directory, DomainsFile)
	wordsFile := filepath.Join(directory, WordsFile)

	domains, derr := domainsFrom(domainsFile)
	if derr != nil {
		return nil, derr
	}

	words, werr := wordsFrom(wordsFile)
	if werr != nil {
		return nil, werr
	}

	return &Artifact{
		domains: domains,
		words:   words,
	}, nil
}

func domainsFrom(filename string) (*set.Set[string], error) {
	s := set.New[string](1024)

	f, ferr := os.Open(filename)
	if ferr != nil {
		return nil, ferr
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "":
			continue
		case line[0] == '#':
			continue
		}
		s.Insert(line)
	}

	return s, scanner.Err()
}

func wordsFrom(filename string) ([]*regexp.Regexp, error) {
	expressions := make([]*regexp.Regexp, 0, 32)

	f, ferr := os.Open(filename)
	if ferr != nil {
		return nil, ferr
	}

	original := make(map[string]string, 32)
	jerr := json.NewDecoder(f).Decode(&original)
	if jerr != nil {
		return nil, jerr
	}

	for _, v := range original {
		reg, cerr := regexp.Compile(v)
		if cerr != nil {
			return nil, cerr
		}
		expressions = append(expressions, reg)
	}

	return expressions, nil
}
