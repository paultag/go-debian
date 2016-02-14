package control_test

import (
	"bytes"
	"strings"
	"testing"

	"pault.ag/go/debian/control"
	// "pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

type TestMarshalStruct struct {
	Foo string
}

type SomeVersionStruct struct {
	Version version.Version
}

type TestParaMarshalStruct struct {
	control.Paragraph
	Foo string
}

func TestExtraMarshal(t *testing.T) {
	el := TestParaMarshalStruct{}

	isok(t, control.Unmarshal(&el, strings.NewReader(`Foo: test
X-A-Test: Foo
`)))

	assert(t, el.Foo == "test")

	writer := bytes.Buffer{}
	isok(t, control.Marshal(&writer, el))
	assert(t, writer.String() == `Foo: test
X-A-Test: Foo

`)
}

func TestBasicMarshal(t *testing.T) {
	testStruct := TestMarshalStruct{Foo: "Hello"}
	writer := bytes.Buffer{}

	err := control.Marshal(&writer, testStruct)
	isok(t, err)

	assert(t, testStruct.Foo == "Hello")
	assert(t, writer.String() == `Foo: Hello

`)
}

func TestExternalMarshal(t *testing.T) {
	testStruct := SomeVersionStruct{}
	isok(t, control.Unmarshal(&testStruct, strings.NewReader(`Version: 1.0-1
`)))
	writer := bytes.Buffer{}

	err := control.Marshal(&writer, testStruct)
	isok(t, err)

	assert(t, writer.String() == `Version: 1.0-1

`)
}

// vim: foldmethod=marker
