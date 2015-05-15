package control

import (
	"bufio"
	"fmt"

	"pault.ag/x/go-debian/dependency"
)

type Control struct {
	Source   SourceParagraph
	Binaries []BinaryParagraph
}

type SourceParagraph struct {
	Paragraph

	Maintainer          string
	Source              string
	BuildDepends        dependency.Dependency
	BuildDependsIndep   dependency.Dependency
	BuildConflicts      dependency.Dependency
	BuildConflictsIndep dependency.Dependency
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

type BinaryParagraph struct {
	Paragraph
	Arch        dependency.Arch
	Description string
	Package     string
}

func ParseControl(reader *bufio.Reader) (ret *Control, err error) {
	ret = &Control{
		Binaries: []BinaryParagraph{},
	}

	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	ret.Source = SourceParagraph{
		Paragraph:  *src,
		Maintainer: src.Values["Maintainer"],
		Source:     src.Values["Source"],

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
