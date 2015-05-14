package control_test

import (
	"bufio"
	"log"
	"strings"
	"testing"

	"pault.ag/x/go-debian/control"
)

/*
 *
 */

func isok(t *testing.T, err error) {
	if err != nil {
		log.Printf("Error! Error is not nil!\n")
		t.FailNow()
	}
}

func notok(t *testing.T, err error) {
	if err == nil {
		log.Printf("Error! Error is  nil!\n")
		t.FailNow()
	}
}

func assert(t *testing.T, expr bool) {
	if !expr {
		log.Printf("Assertion failed!")
		t.FailNow()
	}
}

/*
 *
 */

func TestBasicControlParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
`))
	deb822, err := control.ParseDeb822(reader)
	isok(t, err)

	assert(t, len(deb822.Order) == 1)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")

	reader = bufio.NewReader(strings.NewReader(`Foo: bar`))
	deb822, err = control.ParseDeb822(reader)
	isok(t, err)

	assert(t, len(deb822.Order) == 0)
	assert(t, len(deb822.Order) == len(deb822.Values))
}

func TestMultilineControlParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
Bar-Baz: fnord
 and
 second
 line
 here
`))
	deb822, err := control.ParseDeb822(reader)
	isok(t, err)

	assert(t, len(deb822.Order) == 2)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")
}

func TestFunControlParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
Bar-Baz: fnord:and:this:here
 and:not
 second:here
 line
 here
`))
	deb822, err := control.ParseDeb822(reader)
	isok(t, err)

	assert(t, len(deb822.Order) == 2)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")
}

func TestSingleParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
Bar-Baz: fnord:and:this:here
 and:not
 second:here
 line
 here
Hello: world

But-not: me
`))
	deb822, err := control.ParseDeb822(reader)
	isok(t, err)

	assert(t, len(deb822.Order) == 3)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")

	if _, ok := deb822.Values["But-not"]; ok {
		/* In the case where we *have* this key; let's abort this guy
		 * quickly. */
		t.FailNow()
	}
}
