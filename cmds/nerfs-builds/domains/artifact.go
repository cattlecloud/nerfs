package domains

import (
	"cmp"
	"io"

	"github.com/hashicorp/go-set/v3"
)

// NewArtifact creates an empty artifact in which unwanted domains can be
// tracked for further consumption.
func NewArtifact() *Artifact {
	return &Artifact{
		domains: set.NewTreeSet(cmp.Compare[string]),
	}
}

// Artifact contains a set of domains.
type Artifact struct {
	domains *set.TreeSet[string]
}

// Add each domain to the artifact.
//
// Not thread safe.
func (a *Artifact) Add(domains []string) {
	a.domains.InsertSlice(domains)
}

// Write each domain in order to w.
func (a *Artifact) Write(w io.Writer) error {
	for domain := range a.domains.Items() {
		_, err := io.WriteString(w, domain+"\n")
		if err != nil {
			return err
		}
	}
	return nil
}
