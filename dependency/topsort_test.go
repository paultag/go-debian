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

	"pault.ag/x/go-debian/dependency"
)

func TestTopsortEasy(t *testing.T) {
	foo, err := dependency.Parse("bar")
	isok(t, err)

	bar, err := dependency.Parse("baz, quix")
	isok(t, err)

	arch, err := dependency.ParseArch("amd64")
	isok(t, err)

	order, err := dependency.SortDependencies(map[string][]dependency.Possibility{
		"foo": foo.GetPossibilities(*arch),
		"bar": bar.GetPossibilities(*arch),
	})
	isok(t, err)

	assert(t, len(order) == 2)
	assert(t, order[0] == "foo")
	assert(t, order[1] == "bar")
}

func TestTopsortCycle(t *testing.T) {
	foo, err := dependency.Parse("bar")
	isok(t, err)

	bar, err := dependency.Parse("baz, quix, foo")
	isok(t, err)

	arch, err := dependency.ParseArch("amd64")
	isok(t, err)

	_, err = dependency.SortDependencies(map[string][]dependency.Possibility{
		"foo": foo.GetPossibilities(*arch),
		"bar": bar.GetPossibilities(*arch),
	})
	notok(t, err)
	/* since foo -> bar -> foo; this should raise an error */
}

// vim: foldmethod=marker
