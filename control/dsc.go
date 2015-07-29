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
	"os"
	"path/filepath"
	"strings"

	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"

	"pault.ag/go/topsort"
)

// A DSC is the ecapsulation of a Debian .dsc control file. This contains
// information about the source package, and is general handy.
type DSC struct {
	Paragraph

	Filename string

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

func OrderDSCForBuild(dscs []*DSC, arch dependency.Arch) ([]*DSC, error) {
	sourceMapping := map[string]string{}
	network := topsort.NewNetwork()
	ret := []*DSC{}

	/*
	 * - Create binary -> source mapping.
	 * - Error if two sources provide the same binary
	 * - Create a node for each source
	 * - Create an edge from the source -> source
	 * - return sorted list of dsc files
	 */

	for _, dsc := range dscs {
		for _, binary := range dsc.Binaries {
			sourceMapping[binary] = dsc.Source
		}
		network.AddNode(dsc.Source, dsc)
	}

	for _, dsc := range dscs {
		concreteBuildDepends := dsc.BuildDepends.GetPossibilities(arch)
		for _, relation := range concreteBuildDepends {
			if val, ok := sourceMapping[relation.Name]; ok {
				err := network.AddEdge(val, dsc.Source)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	nodes, err := network.Sort()
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		ret = append(ret, node.Value.(*DSC))
	}

	return ret, nil
}

func ParseDscFile(path string) (ret *DSC, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ret, err = ParseDsc(bufio.NewReader(f), path)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Given a bufio.Reader, produce a DSC struct to encapsulate the
// data contained within.
func ParseDsc(reader *bufio.Reader, path string) (ret *DSC, err error) {

	/* a DSC is a Paragraph, with some stuff. So, let's first take
	 * the bufio.Reader and produce a stock Paragraph. */
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

	arch, err := dependency.ParseArchitectures(src.Values["Architecture"])
	if err != nil {
		return nil, err
	}

	ret = &DSC{
		Paragraph: *src,
		Filename:  path,

		Format: src.Values["Format"],
		Source: src.Values["Source"],

		Architectures:    arch,
		Version:          version,
		Origin:           src.Values["Origin"],
		Maintainer:       src.Values["Maintainer"],
		Homepage:         src.Values["Homepage"],
		StandardsVersion: src.Values["Standards-Version"],
		Binaries:         strings.Split(src.Values["Binary"], ", "),

		Maintainers:  maintainers,
		Uploaders:    uploaders,
		BuildDepends: src.getOptionalDependencyField("Build-Depends"),
	}

	return
}

// vim: foldmethod=marker
