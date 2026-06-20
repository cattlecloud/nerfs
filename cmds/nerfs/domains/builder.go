package domains

import (
	"bytes"
	"path/filepath"
	"time"

	"cattlecloud.net/go/atomicfs"
	"cattlecloud.net/go/ulog"
	"cattlecloud.net/nerfs"
)

var sources = []string{
	// malware / adware
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",

	// gambling sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling-only/hosts",

	// adult sites
	"https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn-only/hosts",
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
