package control_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
)

type AnotherTestStruct struct {
	Value      string   `required:"true"`
	ValueTwo   string   `control:"Value-Two"`
	ValueThree []string `delim:","`
	ValueFour  []string `delim:"|"`
	Depends    dependency.Dependency
}

func TestMarshalRoundTrip(t *testing.T) {
	foo := AnotherTestStruct{}
	foo.Value = "true"
	foo.ValueTwo = "bar"
	foo.ValueThree = []string{"three", "four", "five"}
	foo.ValueFour = []string{"six", "seven", "eight"}
	dep, err := dependency.Parse("foo, bar, baz")
	isok(t, err)
	foo.Depends = *dep

	data := bytes.Buffer{}

	isok(t, control.Marshal(&foo, &data))

	newFoo := AnotherTestStruct{}
	isok(t, control.Unmarshal(&newFoo, &data))

	assert(t, foo.Value == newFoo.Value)
	assert(t, foo.Depends.Relations[0].Possibilities[0].Name == "foo")
	assert(t, newFoo.ValueThree[1] == "four")
	assert(t, newFoo.ValueFour[1] == "seven")
}

func TestMarshalSliceRoundTrip(t *testing.T) {
	foos := []AnotherTestStruct{
		AnotherTestStruct{Value: "Foo", ValueTwo: "FooBar"},
		AnotherTestStruct{Value: "Bar", ValueTwo: "BarFoo"},
	}
	newFoos := []AnotherTestStruct{}

	data := bytes.Buffer{}

	isok(t, control.Marshal(&foos, &data))
	fmt.Printf("'%s'\n", data.String())
	reader := bufio.NewReader(&data)
	isok(t, control.Unmarshal(&newFoos, reader))

	fmt.Printf("%s\n", newFoos)

	assert(t, foos[0].Value == newFoos[0].Value)
	assert(t, foos[1].Value == newFoos[1].Value)
}
