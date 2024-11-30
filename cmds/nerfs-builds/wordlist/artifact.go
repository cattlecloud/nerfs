package wordlist

import (
	"cmp"
	"io"
	"strings"

	"github.com/hashicorp/go-set/v3"
)

type Artifact struct {
	words *set.TreeSet[string]
}

func NewArtifact() *Artifact {
	return &Artifact{
		words: set.NewTreeSet(cmp.Compare[string]),
	}
}

func (a *Artifact) Add(word string) {
	a.words.Insert(word)
	a.words.Insert(word + "s")
	a.words.Insert(strings.ReplaceAll(word, " ", "-"))
	a.words.Insert(strings.ReplaceAll(word, "i", "1"))
	a.words.Insert(strings.ReplaceAll(word, "o", "0"))
	a.words.Insert(strings.ReplaceAll(word, "a", "@"))
	a.words.Insert(strings.ReplaceAll(word, "s", "$"))
}

func (a *Artifact) Write(w io.Writer) error {
	for word := range a.words.Items() {
		_, err := io.WriteString(w, word+"\n")
		if err != nil {
			return err
		}
	}
	return nil
}
