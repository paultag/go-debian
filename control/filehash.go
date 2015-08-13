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
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

func hashFile(path string, algo hash.Hash) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(algo, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(algo.Sum(nil)), nil
}

type DebianFileHash struct {
	// cb136f28a8c971d4299cc68e8fdad93a8ca7daf3 1131 dput-ng_1.9.dsc
	Algorithm string
	Hash      string
	Size      int
	Filename  string
}

func (d *DebianFileHash) Validate() (bool, error) {
	var algo hash.Hash

	stat, err := os.Stat(d.Filename)
	if err != nil {
		return false, err
	}
	if size := stat.Size(); size != int64(d.Size) {
		return false, fmt.Errorf("Error! Size mismatch! %d != %d", size, d.Size)
	}

	switch d.Algorithm {
	case "md5":
		algo = md5.New()
	case "sha1":
		algo = sha1.New()
	case "sha256":
		algo = sha256.New()
	default:
		return false, fmt.Errorf("Unknown algorithm: %s", d.Algorithm)
	}
	fileHash, err := hashFile(d.Filename, algo)
	if err != nil {
		return false, err
	}
	return fileHash == d.Hash, nil
}

// {{{ SHA DebianFileHash (both 1 and 256)

type SHADebianFileHash struct {
	DebianFileHash
}

func (c *SHADebianFileHash) unmarshalControl(algorithm, data string) error {
	var err error
	c.Algorithm = algorithm
	vals := strings.Split(data, " ")
	if len(data) < 4 {
		return fmt.Errorf("Error: Unknown SHA Hash line: '%s'", data)
	}

	c.Hash = vals[0]
	c.Size, err = strconv.Atoi(vals[1])
	if err != nil {
		return err
	}
	c.Filename = vals[2]
	return nil
}

// {{{ SHA1 DebianFileHash

type SHA1DebianFileHash struct{ SHADebianFileHash }

func (c *SHA1DebianFileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("sha1", data)
}

// }}}

// {{{ SHA256 DebianFileHash

type SHA256DebianFileHash struct{ SHADebianFileHash }

func (c *SHA256DebianFileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("sha256", data)
}

// }}}

// }}}

// vim: foldmethod=marker
