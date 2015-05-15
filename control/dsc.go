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

package control

import (
	"bufio"

	"pault.ag/x/go-debian/dependency"
	"pault.ag/x/go-debian/version"
)

type DSC struct {
	Paragraph

	Format           string
	Source           string
	Binaries         []string
	Architectures    []dependency.Arch
	Version          version.Version
	Origin           string
	Maintainer       string
	Maintainers      []string
	Uploaders        []string
	Homepage         string
	StandardsVersion string
	BuildDepends     dependency.Dependency

	/*
		TODO:
			Package-List
			Checksums-Sha1
			Checksums-Sha256
			Files
	*/
}

func ParseDsc(reader *bufio.Reader) (ret *DSC, err error) {

	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	uploaders := splitList(src.Values["Uploaders"])
	maintainers := append(uploaders, src.Values["Maintainer"])
	version, err := version.Parse(src.Values["Version"])
	if err != nil {
		return nil, err
	}

	ret = &DSC{
		Paragraph: *src,

		Format: src.Values["Format"],
		Source: src.Values["Source"],

		// Binaries:
		// Architecturess:

		Version:          *version,
		Origin:           src.Values["Origin"],
		Maintainer:       src.Values["Maintainer"],
		Homepage:         src.Values["Homepage"],
		StandardsVersion: src.Values["Standards-Version"],

		Maintainers:  maintainers,
		Uploaders:    uploaders,
		BuildDepends: src.getOptionalDependencyField("Build-Depends"),
	}

	return
}
