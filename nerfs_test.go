package nerfs

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-set/v3"
	"github.com/shoenig/test/must"
)

const (
	sampleWords = `{
  "poop": "p[oO0]+p",
  "fart": "f[a4]rt"
}`
	sampleDomains = `
example.org
example.com

# some comment
example.xyz
`
)

func writeDomains(t *testing.T, directory string) {
	domainsFile, ferr := os.OpenFile(
		filepath.Join(directory, DomainsFile),
		os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	must.NoError(t, ferr)

	_, werr := io.WriteString(domainsFile, sampleDomains)
	must.NoError(t, werr)
	must.Close(t, domainsFile)
}

func writeWords(t *testing.T, directory string) {
	wordsFile, ferr := os.OpenFile(
		filepath.Join(directory, WordsFile),
		os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	must.NoError(t, ferr)

	_, werr := io.WriteString(wordsFile, sampleWords)
	must.NoError(t, werr)
	must.Close(t, wordsFile)
}

func TestArtifact_Load(t *testing.T) {
	t.Parallel()

	d := t.TempDir()

	writeDomains(t, d)
	writeWords(t, d)

	art, lerr := Load(d)
	must.NoError(t, lerr)
	must.Size(t, 3, art.domains)
	must.Len(t, 2, art.words)
}

func TestSynopsis_matching(t *testing.T) {
	t.Parallel()

	a := &Artifact{
		domains: set.From([]string{"example.com", "example.org", "example.xyz"}),
		words: []*regexp.Regexp{
			regexp.MustCompile(`p[oO0]+p`),
			regexp.MustCompile(`f[a4]rt`),
		},
	}

	t.Run("no match", func(t *testing.T) {
		r := strings.NewReader("this is some text! wow!")
		syn := a.Synopsis(r)
		must.Zero(t, syn.Domains)
		must.Zero(t, syn.Words)
		must.False(t, syn.Any())
	})

	t.Run("block domain", func(t *testing.T) {
		r := strings.NewReader("visit example.org today! wow!")
		syn := a.Synopsis(r)
		must.One(t, syn.Domains)
		must.Zero(t, syn.Words)
		must.True(t, syn.Any())
	})

	t.Run("block word", func(t *testing.T) {
		r := strings.NewReader("take a big p00p! wow!")
		syn := a.Synopsis(r)
		must.Zero(t, syn.Domains)
		must.One(t, syn.Words)
		must.True(t, syn.Any())
	})

	t.Run("block mix", func(t *testing.T) {
		text := `
visit example.org we have pOop and f4rts! example.xyz and example.com too!
`
		r := strings.NewReader(text)
		syn := a.Synopsis(r)
		must.Eq(t, 3, syn.Domains)
		must.Eq(t, 2, syn.Words)
		must.True(t, syn.Any())
	})
}

func TestArtifact_matchDomain(t *testing.T) {
	t.Parallel()

	a := &Artifact{
		domains: set.From([]string{"example.com", "example.xyz"}),
	}

	cases := []struct {
		input string
		exp   bool
	}{
		{
			input: "example.com",
			exp:   true,
		},
		{
			input: "https://example.com",
			exp:   true,
		},
		{
			input: "https://example.com/path",
			exp:   true,
		},
		{
			input: "https://safe.example.com/path",
			exp:   false,
		},
		{
			input: "other",
			exp:   false,
		},
	}

	for _, tc := range cases {
		result := a.matchDomain(tc.input)
		must.Eq(t, tc.exp, result)
	}
}
