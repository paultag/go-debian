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

package control

import (
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

// FileHash {{{

// A FileHash is an entry as found in the Files, Checksum-Sha1, and
// Checksum-Sha256 entry for the .dsc or .changes files.
type FileHash struct {
	// cb136f28a8c971d4299cc68e8fdad93a8ca7daf3 1131 dput-ng_1.9.dsc
	Algorithm string
	Hash      string
	Size      int64
	Filename  string
}

// {{{ Hash File implementations

type boringFileHash struct {
	FileHash
}

func (c *boringFileHash) unmarshalControl(algorithm, data string) error {
	var err error
	c.Algorithm = algorithm
	vals := strings.Fields(data)
	if len(data) < 4 {
		return fmt.Errorf("Error: Unknown Debian Hash line: '%s'", data)
	}

	c.Hash = vals[0]
	c.Size, err = strconv.ParseInt(vals[1], 10, 64)
	if err != nil {
		return err
	}
	c.Filename = vals[2]
	return nil
}

// {{{ MD5 FileHash

type MD5FileHash struct{ boringFileHash }

func (c *MD5FileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("md5", data)
}

// }}}

// {{{ SHA1 FileHash

type SHA1FileHash struct{ boringFileHash }

func (c *SHA1FileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("sha1", data)
}

// }}}

// {{{ SHA256 FileHash

type SHA256FileHash struct{ boringFileHash }

func (c *SHA256FileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("sha256", data)
}

// }}}

// }}}

// }}}

// FileHashValidator {{{

type FileHashValidator struct {
	target hash.Hash
	size   int64
	hash   FileHash
}

// io.Writer interface {{{

func (dh *FileHashValidator) Write(p []byte) (int, error) {
	n, err := dh.target.Write(p)
	dh.size += int64(n)
	return n, err
}

// }}}

// GetHash {{{

func (d FileHash) GetHash() (hash.Hash, error) {
	switch d.Algorithm {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("Unknown algorithm: %s", d.Algorithm)
	}
}

// }}}

// Validate {{{

func (d FileHashValidator) Validate() bool {
	hash := fmt.Sprintf("%x", d.target.Sum(nil))
	return d.hash.Size == d.size && d.hash.Hash == hash
}

// }}}

// FileHash -> FileHashValidator {{{

func (d FileHash) Validator() (*FileHashValidator, error) {
	hasher, err := d.GetHash()
	if err != nil {
		return nil, err
	}

	return &FileHashValidator{
		target: hasher,
		size:   0,
		hash:   d,
	}, nil
}

// }}}

// }}}

type FileHashes []FileHash
type FileHashValidators []*FileHashValidator

func (f FileHashes) Validators() (FileHashValidators, error) {
	ret := FileHashValidators{}
	for _, el := range f {
		validator, err := el.Validator()
		if err != nil {
			return FileHashValidators{}, err
		}
		ret = append(ret, validator)
	}
	return ret, nil
}

func (f FileHashValidators) Writer() io.Writer {
	writers := []io.Writer{}
	for _, el := range f {
		writers = append(writers, el)
	}
	return io.MultiWriter(writers...)
}

func (f FileHashValidators) Validate() bool {
	for _, el := range f {
		if ok := el.Validate(); !ok {
			return false
		}
	}
	return true
}

// vim: foldmethod=marker
