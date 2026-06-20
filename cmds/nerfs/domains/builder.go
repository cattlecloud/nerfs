package domains

import (
	"bytes"
	"path/filepath"
	"time"

	"cattlecloud.net/go/atomicfs"
	"cattlecloud.net/go/nerfs"
	"cattlecloud.net/go/ulog"
)

var sources = []string{
	// malware / adware
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",

	// gambling sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling-only/hosts",

	// adult sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn-only/hosts",

	// fake news sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-only/hosts",

	// social sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/social-only/hosts",

	// more reddit sites
	"https://gist.githubusercontent.com/shoenig/a831965d4dcb36e18034e4703c814f4b/raw/2afcb7642570b36a3ed25ebaab611a772eb3f879/blockreddit.txt",

	// more sites
	"https://gist.githubusercontent.com/shoenig/419256bc1e6938945df70e8bdd260e8e/raw/b5cc28e53b3c6e85d2a6a274dc7e544ce1670679/gistfile1.txt",
}

type Builder struct {
	log *ulog.Log
}

func NewBuilder() *Builder {
	return &Builder{
		log: ulog.New("domains-builder"),
	}
}

func (b *Builder) Build(directory string) error {
	start := time.Now()
	b.log.I.Fmt("starting the build ...")

	art := NewArtifact()

	for _, source := range sources {
		s := NewSource(source)
		if err := s.Get(art); err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)
	if err := art.Write(buf); err != nil {
		return err
	}

	fw := atomicfs.New(atomicfs.Options{
		TmpDirectory: directory,
		TmpExtension: ".temp",
		Mode:         0o644,
	})

	destination := filepath.Join(directory, nerfs.DomainsFile)
	b.log.I.Fmt("writing artifact to %s", destination)
	if err := fw.WriteFile(destination, buf); err != nil {
		return err
	}

	elapsed := time.Since(start)
	b.log.I.Fmt("complete in %v", elapsed)
	return nil
}
