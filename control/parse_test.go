/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package control_test

import (
	"bufio"
	"log"
	"strings"
	"testing"

	"pault.ag/go/debian/control"
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
	deb822, err := control.ParseParagraph(reader)
	isok(t, err)
	assert(t, deb822 != nil)

	assert(t, len(deb822.Order) == 1)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")

	reader = bufio.NewReader(strings.NewReader(`Foo: bar`))
	deb822, err = control.ParseParagraph(reader)
	assert(t, deb822 == nil)
	assert(t, err == nil)
}

func TestMultilineControlParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
Bar-Baz: fnord
 and
 second
 line
 here
`))
	deb822, err := control.ParseParagraph(reader)
	isok(t, err)
	assert(t, deb822 != nil)

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
	deb822, err := control.ParseParagraph(reader)
	isok(t, err)
	assert(t, deb822 != nil)

	assert(t, len(deb822.Order) == 2)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")
}

func TestSingleParse(t *testing.T) {
	// Test Paragraph {{{
	reader := bufio.NewReader(strings.NewReader(`Foo: bar
Bar-Baz: fnord:and:this:here
 and:not
 second:here
 line
 here
Hello: world

But-not: me
`))
	// }}}
	deb822, err := control.ParseParagraph(reader)
	isok(t, err)
	assert(t, deb822 != nil)

	assert(t, len(deb822.Order) == 3)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["Foo"] == "bar")

	if _, ok := deb822.Values["But-not"]; ok {
		/* In the case where we *have* this key; let's abort this guy
		 * quickly. */
		t.FailNow()
	}

	deb822, err = control.ParseParagraph(reader)
	isok(t, err)
	assert(t, deb822 != nil)

	assert(t, len(deb822.Order) == 1)
	assert(t, len(deb822.Order) == len(deb822.Values))
	assert(t, deb822.Values["But-not"] == "me")
}

// func TestSeralize(t *testing.T) {
// 	// Test Paragraph {{{
// 	para := `Foo: bar
// Bar-Baz: fnord:and:this:here
//  and:not
//  second:here
//  line
//  here
// Hello: world
//
// But-not: me
// `
// 	// }}}
// 	reader := bufio.NewReader(strings.NewReader(para))
// 	deb822, err := control.ParseParagraph(reader)
// 	isok(t, err)
//
// 	out := deb822.String()
// 	assert(t, out == `Foo: bar
// Bar-Baz: fnord:and:this:here
//  and:not
//  second:here
//  line
//  here
// Hello: world
// `)
// }

// vim: foldmethod=marker
