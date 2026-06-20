package wordlist

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

type Artifact struct {
	words map[string]*regexp.Regexp
}

func NewArtifact() *Artifact {
	return &Artifact{
		words: make(map[string]*regexp.Regexp, 32),
	}
}

func (a *Artifact) Add(line string) {
	// given line is directly from words.txt which is in the form
	// <word> -> <regexp>

	tokens := strings.Split(line, "->")
	if len(tokens) != 2 {
		return
	}

	original := strings.TrimSpace(tokens[0])
	regex := strings.TrimSpace(tokens[1])

	re := regexp.MustCompile(regex)
	a.words[original] = re
}

func (a *Artifact) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(a.words)
}
