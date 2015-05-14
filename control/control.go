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
}

func ParseControl(reader *bufio.Reader) (ret *Control, err error) {
	ret = &Control{
		Binaries: []ControlBinary{},
	}

	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	ret.Source = ControlSource{*src}

	for {
		para, err := ParseParagraph(reader)
		if err != nil {
			return nil, err
		}
		if para == nil {
			break
		}

		ret.Binaries = append(ret.Binaries, ControlBinary{*para})
	}
	return
}
