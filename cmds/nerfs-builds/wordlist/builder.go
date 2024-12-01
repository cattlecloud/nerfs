package wordlist

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
	"time"

	"noxide.lol/go/atomicfs"
	"noxide.lol/go/ulog"
)

type Builder struct {
	log *ulog.Log
}

func NewBuilder() *Builder {
	return &Builder{
		log: ulog.New("wordlist-builder"),
	}
}

func (b *Builder) Build(destination string) error {
	start := time.Now()
	b.log.I.Fmt("starting the build ...")

	art := NewArtifact()

	reader := strings.NewReader(Source)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		art.Add(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := art.Write(buf); err != nil {
		return err
	}

	directory := filepath.Dir(destination)
	fw := atomicfs.New(atomicfs.Options{
		TmpDirectory: directory,
		TmpExtension: ".temp",
		Mode:         0o644,
	})

	b.log.I.Fmt("writing artifact to %s", destination)
	if err := fw.WriteFile(destination, buf); err != nil {
		return err
	}

	elapsed := time.Since(start)
	b.log.I.Fmt("complete in %v", elapsed)
	return nil
}
