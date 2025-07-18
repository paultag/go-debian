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

func TestArchBasics(t *testing.T) {
	arch, err := dependency.ParseArch("amd64")
	isok(t, err)
	assert(t, arch.CPU == "amd64")
	assert(t, arch.ABI == "base")
	assert(t, arch.Libc == "gnu")
	assert(t, arch.OS == "linux")
}

func TestArchCompareX32(t *testing.T) {
	archX32, _ := dependency.ParseArch("x32")
	archAnyAmd64, _ := dependency.ParseArch("any-amd64")
	assert(t, archX32.Is(archAnyAmd64))
}

/*
 */
func TestArchCompareBasics(t *testing.T) {
	arch, err := dependency.ParseArch("amd64")
	isok(t, err)

	equivs := []string{
		"base-gnu-linux-amd64",
		"linux-any",
		"any",
		"base-gnu-linux-any",
	}

	for _, el := range equivs {
		other, err := dependency.ParseArch(el)
		isok(t, err)
		assert(t, arch.Is(other))
		assert(t, other.Is(arch))
	}

	unequivs := []string{
		"base-gnu-linux-all",
		"all",

		"base-gnuu-linux-amd64",
		"base-gnu-linuxx-amd64",
		"base-gnu-linux-amd644",
	}

	for _, el := range unequivs {
		other, err := dependency.ParseArch(el)
		isok(t, err)

		assert(t, !arch.Is(other))
		assert(t, !other.Is(arch))
	}
}

func TestArchCompareAllAny(t *testing.T) {
	all, err := dependency.ParseArch("all")
	if err != nil {
		t.Fatal(err)
	}

	any, err := dependency.ParseArch("any")
	if err != nil {
		t.Fatal(err)
	}

	if all.Is(any) {
		t.Fatalf("arch:all unexpectedly is arch:any")
	}

	if any.Is(all) {
		t.Fatalf("arch:all unexpectedly is arch:any")
	}
}

/*
 */
func TestArchSetCompare(t *testing.T) {
	dep, err := dependency.Parse("foo [amd64], bar [!sparc]")
	isok(t, err)

	iAm, err := dependency.ParseArch("amd64")
	isok(t, err)

	fooArch := dep.Relations[0].Possibilities[0].Architectures
	barArch := dep.Relations[1].Possibilities[0].Architectures

	assert(t, fooArch.Matches(iAm))
	assert(t, barArch.Matches(iAm))

	iAmNot, err := dependency.ParseArch("armhf")
	isok(t, err)

	assert(t, !fooArch.Matches(iAmNot))
	assert(t, barArch.Matches(iAmNot))
}

// vim: foldmethod=marker
