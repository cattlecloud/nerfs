package wordlist

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/shoenig/test/util"
)

const (
	// numWords must match the number of expressions in words.txt
	numWords = 28
)

func TestBuilder_Build(t *testing.T) {
	t.Parallel()

	dest := util.TempFile(t)

	b := NewBuilder()
	err := b.Build(dest)
	must.NoError(t, err)

	f, ferr := os.Open(dest)
	must.NoError(t, ferr)

	m := make(map[string]*regexp.Regexp)
	jerr := json.NewDecoder(f).Decode(&m)
	must.NoError(t, jerr)
	must.MapLen(t, numWords, m)

	t.Run("spot checks", func(t *testing.T) {
		must.RegexMatch(t, m["wop"], "wop")
		must.RegexMatch(t, m["coom"], "c00m$")
		must.RegexMatch(t, m["camel jockey"], "c@m3l-jock3ys")
		must.RegexMatch(t, m["dago"], "d@g0$")

		// non-matching due to word boundary
		matches := m["wop"].MatchString("wopping")
		must.False(t, matches)
	})
}
