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
	"testing"

	"pault.ag/go/debian/control"
)

/*
 *
 */

func TestFilehash(t *testing.T) {
	hash := control.FileHash{
		Algorithm: "sha256",
		Hash:      "1eb9d182f09f015790915a9875005b70328c96ea43010c04b3b2b3a6ff37d2f4",
		Size:      39,
		Filename:  "<nowhere in particular>",
	}
	validator, err := hash.Validator()
	isok(t, err)

	// `nobody inspects the spammish repetition`

	assert(t, !validator.Validate())

	c, err := validator.Write([]byte("nobody inspects the spammish repetition"))
	isok(t, err)
	assert(t, c == 39)

	assert(t, validator.Validate())

}

// vim: foldmethod=marker
