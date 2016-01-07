package control_test

import (
	"bytes"
	"testing"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
)

type AnotherTestStruct struct {
	Value      string `required:"true"`
	ValueTwo   string `control:"Value-Two"`
	ValueThree []string
	Depends    dependency.Dependency
}

func TestMarshalRoundTrip(t *testing.T) {
	foo := AnotherTestStruct{}
	foo.Value = "true"
	foo.ValueTwo = "bar"
	foo.ValueThree = []string{"three"}

	data := bytes.Buffer{}

	isok(t, control.Marshal(&foo, &data))

	newFoo := AnotherTestStruct{}
	isok(t, control.Unmarshal(&newFoo, &data))

	assert(t, foo.Value == newFoo.Value)
}
