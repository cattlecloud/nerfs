package nerfs

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-set/v3"
)

const (
	DomainsFile = "domains.txt"
	WordsFile   = "wordlist.txt"
)

// An Artifact contains the in-memory optimized form of the domain and word
// block-list files created by the nerfs tool.
//
// Create an Artifact by calling Load() with the directory of the compiled
// artifacts.
type Artifact struct {
	domains *set.Set[string]
	words   []*regexp.Regexp
}

// A Synopsis contains the result of analyzing some text.
type Synopsis struct {
	// Domains is the count of matching domains on the block list.
	Domains int

	// Words is the count of matching words on the block list.
	Words int
}

// Any returns whether any domains or words appear on the block lists.
func (s *Synopsis) Any() bool {
	return s.Domains > 0 || s.Words > 0
}

// Synopsis analyzes the content in the given reader and identifies any domains
// or words on one of the block lists.
//
// Note that this should be used carefully - scanning large text consumes lots
// of CPU due to regular expression matching against every token. Many use cases
// need only scan against blocked domains which is must faster.
func (a *Artifact) Synopsis(r io.Reader) *Synopsis {
	syn := new(Synopsis)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		switch {
		case a.matchDomain(word):
			syn.Domains++
		case a.matchWord(word):
			syn.Words++
		}
	}
	return syn
}

func (a *Artifact) matchDomain(s string) bool {
	s, _ = strings.CutPrefix(s, "https://")
	s, _ = strings.CutPrefix(s, "http://")
	if i := strings.IndexAny(s, "/?#"); i >= 0 {
		s = s[:i]
	}
	return a.domains.Contains(s)
}

func (a *Artifact) matchWord(s string) bool {
	for _, reg := range a.words {
		if reg.MatchString(s) {
			return true
		}
	}
	return false
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
