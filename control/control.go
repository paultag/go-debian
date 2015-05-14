package control

import (
	"bufio"
	"fmt"

	"pault.ag/x/go-debian/dependency"
)

type Control struct {
	Source   ControlSource
	Binaries []ControlBinary
}

type ControlSource struct {
	Paragraph
	Maintainer string
	Source     string
}

func (control *ControlSource) GetBuildDepends() (*dependency.Depedency, error) {
	return control.getDependencyField("Build-Depends")
}

func (control *ControlSource) GetBuildDependsIndep() (*dependency.Depedency, error) {
	return control.getDependencyField("Build-Depends-Indep")
}

func (para *Paragraph) getDependencyField(field string) (*dependency.Depedency, error) {
	if val, ok := para.Values[field]; ok {
		return dependency.Parse(val)
	}
	return nil, fmt.Errorf("Field `%s' Missing", field)
}

type ControlBinary struct {
	Paragraph
	Arch        dependency.Arch
	Description string
	Package     string
}

func ParseControl(reader *bufio.Reader) (ret *Control, err error) {
	ret = &Control{
		Binaries: []ControlBinary{},
	}

	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	ret.Source = ControlSource{
		Paragraph:  *src,
		Maintainer: src.Values["Maintainer"],
		Source:     src.Values["Source"],
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

		ret.Binaries = append(ret.Binaries, ControlBinary{
			Paragraph: *para,
			Arch:      *arch,

			Description: para.Values["Description"],
			Package:     para.Values["Package"],
		})
	}
	return
}
