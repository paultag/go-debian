package deb

import (
	"fmt"
	"io"
	"strings"

	"archive/tar"

	"compress/bzip2"
	"compress/gzip"
	"xi2.org/x/xz"
)

type compressionReader func(io.Reader) (io.Reader, error)

func gzipNewReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

func xzNewReader(r io.Reader) (io.Reader, error) {
	return xz.NewReader(r, 0)
}

func bzipNewReader(r io.Reader) (io.Reader, error) {
	return bzip2.NewReader(r), nil
}

var knownCompressionAlgorithms = map[string]compressionReader{
	".tar.gz":  gzipNewReader,
	".tar.bz2": bzipNewReader,
	".tar.xz":  xzNewReader,
}

func (e *DebEntry) IsTarfile() bool {
	return e.getCompressionReader() != nil
}

func (e *DebEntry) getCompressionReader() *compressionReader {
	for key, decompressor := range knownCompressionAlgorithms {
		if strings.HasSuffix(e.Name, key) {
			return &decompressor
		}
	}
	return nil
}

func (e *DebEntry) Tarfile() (*tar.Reader, error) {
	decompressor := e.getCompressionReader()
	if decompressor == nil {
		return nil, fmt.Errorf("%s appears to not be a tarfile", e.Name)
	}
	reader, err := (*decompressor)(e.Data)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(reader), nil
}
