/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package deb

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"archive/tar"

	"compress/bzip2"
	"compress/gzip"

	"github.com/lxq/lzma"
	"github.com/xi2/xz"
)

// known compression types {{{

type DecompressorFunc func(io.Reader) (io.Reader, error)

func gzipNewReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

func xzNewReader(r io.Reader) (io.Reader, error) {
	return xz.NewReader(r, 0)
}

func lzmaNewReader(r io.Reader) (io.Reader, error) {
	return lzma.NewReader(r), nil
}

func bzipNewReader(r io.Reader) (io.Reader, error) {
	return bzip2.NewReader(r), nil
}

// For the authoritative list of supported file formats, see
// https://manpages.debian.org/unstable/dpkg-dev/deb.5
var knownCompressionAlgorithms = map[string]DecompressorFunc{
	".gz":   gzipNewReader,
	".bz2":  bzipNewReader,
	".xz":   xzNewReader,
	".lzma": lzmaNewReader,
}

// DecompressorFn returns a decompressing reader for the specified reader and its
// corresponding file extension ext.
func DecompressorFor(ext string) DecompressorFunc {
	if fn, ok := knownCompressionAlgorithms[ext]; ok {
		return fn
	}
	return func(r io.Reader) (io.Reader, error) { return r, nil } // uncompressed file or unknown compression scheme
}

// }}}

// IsTarfile {{{

// Check to see if the given ArEntry is, in fact, a Tarfile. This method
// will return `true` for `control.tar.*` and `data.tar.*` files.
//
// This will return `false` for the `debian-binary` file. If this method
// returns `true`, the `.Tarfile()` method will be around to give you a
// tar.Reader back.
func (e *ArEntry) IsTarfile() bool {
	ext := filepath.Ext(e.Name)
	return filepath.Ext(strings.TrimSuffix(e.Name, ext)) == ".tar"
}

// }}}

// Tarfile {{{

// `.Tarfile()` will return a `tar.Reader` created from the ArEntry member
// to allow further inspection of the contents of the `.deb`.
func (e *ArEntry) Tarfile() (*tar.Reader, error) {
	if !e.IsTarfile() {
		return nil, fmt.Errorf("%s appears to not be a tarfile", e.Name)
	}
	reader, err := DecompressorFor(filepath.Ext(e.Name))(e.Data)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(reader), nil
}

// }}}

// vim: foldmethod=marker
