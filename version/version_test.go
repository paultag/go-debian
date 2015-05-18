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

package version_test

import (
	"testing"

	"pault.ag/x/go-debian/version"
)

/*
 *
 */

func isok(t *testing.T, err error) {
	if err != nil {
		log.Printf("Error! Error is not nil!\n")
		t.FailNow()
	}
}

func notok(t *testing.T, err error) {
	if err == nil {
		log.Printf("Error! Error is  nil!\n")
		t.FailNow()
	}
}

func assert(t *testing.T, expr bool) {
	if !expr {
		log.Printf("Assertion failed!")
		t.FailNow()
	}
}

/*
 *
 */

func TestSingleParse(t *testing.T) {
	ver, err := version.Parse("1:1.0-1")
	isok(t, err)

	assert(t, !ver.Native)
	assert(t, ver.Epoch == 1)
	assert(t, ver.Version == "1.0")
	assert(t, ver.Revision == "1")
}

func TestNativeParse(t *testing.T) {
	ver, err := version.Parse("1:1.0")
	isok(t, err)

	assert(t, ver.Native)
	assert(t, ver.Epoch == 1)
	assert(t, ver.Version == "1.0")
}

func TestNoEpoch(t *testing.T) {
	ver, err := version.Parse("1.0-1")
	isok(t, err)
	assert(t, !ver.Native)
	assert(t, ver.Epoch == 0)
	assert(t, ver.Version == "1.0")
	assert(t, ver.Revision == "1")

	ver, err = version.Parse("1.0")
	isok(t, err)
	assert(t, ver.Native)
	assert(t, ver.Epoch == 0)
	assert(t, ver.Version == "1.0")
}

func TestComplexParse(t *testing.T) {
	ver, err := version.Parse("1:1:1.0-1-1") // Yes, this is valid.
	isok(t, err)

	assert(t, !ver.Native)
	assert(t, ver.Epoch == 1)
	assert(t, ver.Version == "1:1.0-1")
	assert(t, ver.Revision == "1")
}

// vim: foldmethod=marker
