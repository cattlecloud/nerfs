package domains

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shoenig/test/must"
)

func TestSource_Get(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "0.0.0.0 1.example.com\n")
		_, _ = io.WriteString(w, "0.0.0.0 2.example.com\n")
		_, _ = io.WriteString(w, "0.0.0.0 3.example.com\n")
	}))
	s := NewSource(ts.URL)

	art := NewArtifact()

	err := s.Get(art)
	must.NoError(t, err)

	b := new(bytes.Buffer)
	err = art.Write(b)
	must.NoError(t, err)

	result := b.String()
	must.StrContains(t, result, "1.example.com")
	must.StrContains(t, result, "2.example.com")
	must.StrContains(t, result, "3.example.com")
}
