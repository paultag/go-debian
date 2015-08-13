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
	"strconv"
	"strings"
)

// {{{ Changes File Hash line struct helpers (Files, SHA1, SHA256)

type DebianFileHash struct {
	// cb136f28a8c971d4299cc68e8fdad93a8ca7daf3 1131 dput-ng_1.9.dsc
	Algorithm string
	Hash      string
	Size      int
	Filename  string
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
// }}}

// vim: foldmethod=marker
