package nerfs

import (
	"io"
	"os"
	"path/filepath"
	"testing"

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
