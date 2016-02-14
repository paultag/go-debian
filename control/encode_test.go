package control_test

import (
	"bytes"
	"testing"

	"pault.ag/go/debian/control"
	// "pault.ag/go/debian/dependency"
	// "pault.ag/go/debian/version"
)

type TestMarshalStruct struct {
	Foo string
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
