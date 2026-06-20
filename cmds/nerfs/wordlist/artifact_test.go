package wordlist

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/shoenig/test/must"
)

func TestArtifact_Add(t *testing.T) {
	t.Parallel()

	art := NewArtifact()

	line := `camel jockey     ->  (?i)\bc[a@]m[e3]l[\s-]*j[o0]ck[e3]y[s$]?\b`
	art.Add(line)

	buf := new(bytes.Buffer)
	err := art.Write(buf)
	must.NoError(t, err)

	// ensure our regexp still works
	m := make(map[string]*regexp.Regexp)
	jerr := json.Unmarshal(buf.Bytes(), &m)
	must.NoError(t, jerr)

	re := m["camel jockey"]
	must.NotNil(t, re)
	must.RegexMatch(t, re, "c@m3l-j0ck3y")
}
