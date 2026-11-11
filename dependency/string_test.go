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

func TestArchString(t *testing.T) {
	equivs := map[string]string{
		"all":                  "all",
		"any":                  "all",
		"amd64":                "amd64",
		"base-gnu-linux-amd64": "amd64",
		"foo-bsd-windows-i386": "foo-bsd-windows-i386",
	}

	for _, el := range equivs {
		arch, err := dependency.ParseArch(el)

		isok(t, err)
		assert(t, arch.String() == equivs[el])
	}
}

func TestSubstvarString(t *testing.T) {
	// A substvar is written as ${name}. Writing it out without the braces
	// turned it into a package name with an architecture qualifier, so a
	// dependency field read from a control file and written back did not
	// parse as the same thing.
	dep, err := dependency.Parse("${misc:Depends}, foo (>= 1.0)")
	isok(t, err)

	assert(t, dep.String() == "${misc:Depends}, foo (>= 1.0)")

	rtDep, err := dependency.Parse(dep.String())
	isok(t, err)

	assert(t, len(rtDep.Relations) == 2)
	assert(t, rtDep.Relations[0].Possibilities[0].Name == "misc:Depends")
	assert(t, rtDep.Relations[0].Possibilities[0].Substvar)
	assert(t, !rtDep.Relations[1].Possibilities[0].Substvar)
}

func TestSubstvarStringRoundTrip(t *testing.T) {
	for _, input := range []string{
		"${shlibs:Depends}",
		"${misc:Depends}, ${shlibs:Depends}",
		"foo, ${misc:Depends}",
		"${a} | ${b}",
		"${misc:Depends}, foo (>= 1.0) | bar",
	} {
		dep, err := dependency.Parse(input)
		isok(t, err)

		assert(t, dep.String() == input)

		rtDep, err := dependency.Parse(dep.String())
		isok(t, err)

		assert(t, len(rtDep.Relations) == len(dep.Relations))
		for i, relation := range dep.Relations {
			assert(t, len(rtDep.Relations[i].Possibilities) == len(relation.Possibilities))
			for j, possi := range relation.Possibilities {
				rtPossi := rtDep.Relations[i].Possibilities[j]
				assert(t, rtPossi.Name == possi.Name)
				assert(t, rtPossi.Substvar == possi.Substvar)
			}
		}
	}
}

// vim: foldmethod=marker
