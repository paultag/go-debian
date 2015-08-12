package control_test

import (
	"strings"
	"testing"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

type TestStruct struct {
	Value      string
	ValueTwo   string `control:"Value-Two"`
	ValueThree []string
	Depends    dependency.Dependency
	Version    version.Version
	Arch       dependency.Arch
	Fnord      struct {
		FooBar string `control:"Fnord-Foo-Bar"`
	}
}

func TestBasicDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Foo-Bar: baz
`)))
	assert(t, foo.Value == "foo")
}

func TestTagDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Value-Two: baz
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.ValueTwo == "baz")
}

func TestDependsDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Depends: foo, bar
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.Depends.Relations[0].Possibilities[0].Name == "foo")

	/* Actually invalid below */
	notok(t, control.Decode(&foo, strings.NewReader(`Depends: foo (>= 1.0) (<= 1.0)
`)))
}

func TestVersionDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Version: 1.0-1
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.Version.Revision == "1")
}

func TestArchDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Arch: amd64
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.Arch.CPU == "amd64")
}

func TestNestedDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
Fnord-Foo-Bar: Thing
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.Fnord.FooBar == "Thing")
}

func TestListDecode(t *testing.T) {
	foo := TestStruct{}
	isok(t, control.Decode(&foo, strings.NewReader(`Value: foo
ValueThree: foo bar baz
`)))
	assert(t, foo.Value == "foo")
	assert(t, foo.ValueThree[0] == "foo")
}
