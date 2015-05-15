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
	"fmt"
	"strings"

	"pault.ag/x/go-debian/dependency"
)

type Control struct {
	Source   SourceParagraph
	Binaries []BinaryParagraph
}

type SourceParagraph struct {
	Paragraph

	Maintainer          string
	Maintainers         []string
	Uploaders           []string
	Source              string
	BuildDepends        dependency.Dependency
	BuildDependsIndep   dependency.Dependency
	BuildConflicts      dependency.Dependency
	BuildConflictsIndep dependency.Dependency
}

type BinaryParagraph struct {
	Paragraph
	Arch        dependency.Arch
	Description string
	Package     string
}

func (para *Paragraph) getDependencyField(field string) (*dependency.Dependency, error) {
	if val, ok := para.Values[field]; ok {
		return dependency.Parse(val)
	}
	return nil, fmt.Errorf("Field `%s' Missing", field)
}

func (para *Paragraph) getOptionalDependencyField(field string) dependency.Dependency {
	val := para.Values[field]
	dep, err := dependency.Parse(val)
	if err != nil {
		return dependency.Dependency{}
	}
	return *dep
}

func splitList(names string) (ret []string) {
	for _, el := range strings.Split(names, ",") {
		el := strings.Trim("\n\r\t ", el)
		if el != "" {
			ret = append(ret, el)
		}
	}
	return ret
}

func ParseControl(reader *bufio.Reader) (ret *Control, err error) {
	ret = &Control{
		Binaries: []BinaryParagraph{},
	}

	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	uploaders := splitList(src.Values["Uploaders"])
	maintainers := append(uploaders, src.Values["Maintainer"])

	ret.Source = SourceParagraph{
		Paragraph:   *src,
		Maintainer:  src.Values["Maintainer"],
		Maintainers: maintainers,
		Uploaders:   uploaders,
		Source:      src.Values["Source"],

		BuildDepends:        src.getOptionalDependencyField("Build-Depends"),
		BuildDependsIndep:   src.getOptionalDependencyField("Build-Depends-Indep"),
		BuildConflicts:      src.getOptionalDependencyField("Build-Conflicts"),
		BuildConflictsIndep: src.getOptionalDependencyField("Build-Conflicts-Indep"),
	}

	for {
		para, err := ParseParagraph(reader)
		if err != nil {
			return nil, err
		}
		if para == nil {
			break
		}

		arch, err := dependency.ParseArch(para.Values["Architecture"])
		if err != nil {
			return nil, err
		}

		ret.Binaries = append(ret.Binaries, BinaryParagraph{
			Paragraph: *para,
			Arch:      *arch,

			Description: para.Values["Description"],
			Package:     para.Values["Package"],
		})
	}
	return
}
