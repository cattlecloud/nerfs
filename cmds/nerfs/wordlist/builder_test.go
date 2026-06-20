package wordlist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"cattlecloud.net/go/nerfs"
	"github.com/shoenig/test/must"
)

const (
	// numWords must match the number of expressions in words.txt
	numWords = 14
)

func TestBuilder_Build(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	b := NewBuilder()
	err := b.Build(dir)
	must.NoError(t, err)

	f, ferr := os.Open(filepath.Join(dir, nerfs.WordsFile))
	must.NoError(t, ferr)

	m := make(map[string]*regexp.Regexp)
	jerr := json.NewDecoder(f).Decode(&m)
	must.NoError(t, jerr)
	must.MapLen(t, numWords, m)

	t.Run("spot checks", func(t *testing.T) {
		must.RegexMatch(t, m["wop"], "wop")
		must.RegexMatch(t, m["dago"], "d@g0$")

		// non-matching due to word boundary
		matches := m["wop"].MatchString("wopping")
		must.False(t, matches)
	})
}
