/* Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
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
 * THE SOFTWARE. */

package control_test

import (
	"bufio"
	"strings"
	"testing"

	"pault.ag/x/go-debian/control"
)

/*
 *
 */

func TestDSCParse(t *testing.T) {
	// Test DSC {{{
	reader := bufio.NewReader(strings.NewReader(`Format: 3.0 (quilt)
Source: fbautostart
Binary: fbautostart
Architecture: any
Version: 2.718281828-1
Maintainer: Paul Tagliamonte <paultag@ubuntu.com>
Homepage: https://launchpad.net/fbautostart
Standards-Version: 3.9.3
Vcs-Browser: http://git.debian.org/?p=collab-maint/fbautostart.git
Vcs-Git: git://git.debian.org/collab-maint/fbautostart.git
Build-Depends: debhelper (>= 9)
Package-List:
 fbautostart deb misc optional arch=any
Checksums-Sha1:
 bc36310c15edc9acf48f0a1daf548bcc6f861372 92748 fbautostart_2.718281828.orig.tar.gz
 eaed7f053dce48d4ad4e442bbb0da73ea1181a26 2356 fbautostart_2.718281828-1.debian.tar.xz
Checksums-Sha256:
 bb2fdfd4a38505905222ee02d8236a594bdf6eaefca23462294cacda631745c1 92748 fbautostart_2.718281828.orig.tar.gz
 f7186d1bebde403527b5b3fd80406decaaf295366206667d5b402da962f0b772 2356 fbautostart_2.718281828-1.debian.tar.xz
Files:
 06495f9b23b1c9b1bf35c2346cb48f63 92748 fbautostart_2.718281828.orig.tar.gz
 f58c0e0bf4d56461e776232484c07301 2356 fbautostart_2.718281828-1.debian.tar.xz
`))
	// }}}
	c, err := control.ParseDsc(reader)
	isok(t, err)
	assert(t, c != nil)

	assert(t, c.Format == "3.0 (quilt)")
	assert(t, c.Source == "fbautostart")
	assert(t, len(c.Maintainers) == 1)
	assert(t, c.Maintainers[0] == "Paul Tagliamonte <paultag@ubuntu.com>")
	assert(t, c.Maintainer == "Paul Tagliamonte <paultag@ubuntu.com>")

	assert(t, c.Version.Version == "2.718281828")
	assert(t, c.Version.Revision == "1")

	assert(t, c.StandardsVersion == "3.9.3")
	assert(t, c.Homepage == "https://launchpad.net/fbautostart")
}

// vim: foldmethod=marker
