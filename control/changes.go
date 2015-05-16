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
	"bufio"

	"pault.ag/x/go-debian/dependency"
	"pault.ag/x/go-debian/version"
)

type Changes struct {
	Paragraph

	Format        string
	Source        string
	Binaries      []string
	Architectures []dependency.Arch
	Version       version.Version
	Origin        string
	Distribution  string
	Urgency       string
	Maintainer    string
	ChangedBy     string
	Closes        []string
	Changes       string

	/*
		TODO:
			Checksums-Sha1
			Checksums-Sha256
			Files
	*/
}

func ParseChanges(reader *bufio.Reader) (ret *Changes, err error) {

	/* a Changes is a Paragraph, with some stuff. So, let's first take
	 * the bufio.Reader and produce a stock Paragraph. */
	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	arch, err := dependency.ParseArchitectures(src.Values["Architecture"])
	if err != nil {
		return nil, err
	}

	version, err := version.Parse(src.Values["Version"])
	if err != nil {
		return nil, err
	}

	ret = &Changes{
		Paragraph: *src,

		Format:     src.Values["Format"],
		Source:     src.Values["Source"],
		Origin:     src.Values["Origin"],
		Urgency:    src.Values["Urgency"],
		Maintainer: src.Values["Maintainer"],
		ChangedBy:  src.Values["Changed-By"],
		Changes:    src.Values["Changes"],

		Architectures: arch,
		Version:       *version,
		Distribution:  src.Values["Distribution"],
		/*
			Binaries      []string
			Closes        []string
		*/
	}

	return
}

// vim: foldmethod=marker
