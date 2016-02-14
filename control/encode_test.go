package control_test

import (
	"bytes"
	"strings"
	"testing"

	"pault.ag/go/debian/control"
	// "pault.ag/go/debian/dependency"
	// "pault.ag/go/debian/version"
)

type TestMarshalStruct struct {
	Foo string
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

// vim: foldmethod=marker
