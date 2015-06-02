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

package dependency_test

import (
	"testing"

	"pault.ag/go/debian/dependency"
)

/*
 *
 */

func TestSliceParse(t *testing.T) {
	dep, err := dependency.Parse("foo, bar | baz")
	isok(t, err)
	arch, err := dependency.ParseArch("amd64")
	isok(t, err)

	els := dep.GetPossibilities(*arch)
	assert(t, len(els) == 2)

	assert(t, els[0].Name == "foo")
	assert(t, els[1].Name == "bar")
}

func TestArchSliceParse(t *testing.T) {
	dep, err := dependency.Parse("foo, bar [sparc] | baz")
	isok(t, err)
	arch, err := dependency.ParseArch("amd64")
	isok(t, err)

	els := dep.GetPossibilities(*arch)
	assert(t, len(els) == 2)

	assert(t, els[0].Name == "foo")
	assert(t, els[1].Name == "baz")
}

func TestSliceAllParse(t *testing.T) {
	dep, err := dependency.Parse("foo, bar | baz")
	isok(t, err)

	els := dep.GetAllPossibilities()
	assert(t, len(els) == 3)

	assert(t, els[0].Name == "foo")
	assert(t, els[1].Name == "bar")
	assert(t, els[2].Name == "baz")
}

func TestSliceSubParse(t *testing.T) {
	dep, err := dependency.Parse("${foo:Depends}, foo, bar | baz, ${bar:Depends}")
	isok(t, err)

	els := dep.GetAllPossibilities()
	assert(t, len(els) == 3)

	assert(t, els[0].Name == "foo")
	assert(t, els[1].Name == "bar")
	assert(t, els[2].Name == "baz")

	els = dep.GetSubstvars()
	assert(t, len(els) == 2)

	assert(t, els[0].Name == "foo:Depends")
	assert(t, els[1].Name == "bar:Depends")
}

// vim: foldmethod=marker
