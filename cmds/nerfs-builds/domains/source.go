package domains

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"time"

	"cattlecloud.net/go/scope"
	"cattlecloud.net/go/ulog"
)

// Source contains a URI from which we can acquire a set of unwanted domains.
type Source struct {
	url    string
	log    *ulog.Log
	client *http.Client
}

func NewSource(url string) *Source {
	return &Source{
		url:    url,
		client: &http.Client{Timeout: 1 * time.Minute},
		log:    ulog.New("source"),
	}
}

func (s *Source) Get(a *Artifact) error {
	ctx, cancel := scope.TTL(30 * time.Second)
	defer cancel()

	request, rerr := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if rerr != nil {
		s.log.E.Fmt("unable to create source request: %v", rerr)
		return rerr
	}
	request.Header.Set("User-Agent", "nerfs/v0")

	response, derr := s.client.Do(request)
	if derr != nil {
		s.log.E.Fmt("unable to make http request: %v", derr)
		return derr
	}
	defer func() { _ = response.Body.Close() }()

	return s.load(response.Body, a)
}

func (s *Source) load(r io.Reader, a *Artifact) error {
	scanner := bufio.NewScanner(r)

	results := make([]string, 0, 1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		chop := strings.TrimPrefix(line, "0.0.0.0")
		clean := strings.TrimSpace(chop)
		results = append(results, clean)
	}

	a.Add(results)

	return scanner.Err()
}
