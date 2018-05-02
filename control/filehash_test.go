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

package control_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/paultag/go-debian/control"
)

func TestCrashWithoutDelim(t *testing.T) {
	type release struct {
		MD5Sum []control.MD5FileHash
		SHA1   []control.SHA1FileHash
		SHA256 []control.SHA256FileHash
	}
	// Test Paragraph {{{
	reader := bufio.NewReader(strings.NewReader(`MD5Sum:
 b096313254606a2f984f6fc38e4769c8  1219120 contrib/Contents-amd64
SHA256:
 0f6dc8c4865a433d321a3af0d4812030fdb14231cd1f93bd7570f1518a46faad  1219120 contrib/Contents-amd64
`))
	// }}}
	var r release
	if err := control.Unmarshal(&r, reader); err == nil {
		t.Errorf("control.Unmarshal unexpectedly succeeded on struct without delim")
	}
}
