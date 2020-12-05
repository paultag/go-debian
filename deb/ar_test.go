package deb_test

import (
	"io"
	"io/ioutil"
	"log"
	"fmt"
	"testing"
	"runtime"
	"pault.ag/go/debian/deb"
)


/* ar is binary format, exact number of spaces is required here */
var foobar_ar = "!<arch>\n" +
	"foo.txt/        0           0     0     644     7         `\n" +
	"FooBar\n\n" +
	"bar.txt         1559837289  1000  1000  100644  21        `\n" +
	"FooBar\nFooBar\nFooBar\n\n"

type wholeReader struct {
	Content []byte
	Pos int
}

func (r *wholeReader) Read(p []byte) (n int, err error) {
	if r.Pos >= len(r.Content) {
		return 0, io.EOF
	}
	n = copy(p, r.Content[r.Pos:])
	r.Pos += n
	return n, nil
}

type byteReader struct {
	Content []byte
	Pos int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.Pos >= len(r.Content) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = r.Content[r.Pos]
	r.Pos += 1
	return 1, nil
}

func failCaller2(t *testing.T, msg string, a ...interface{}) {
	outmsg := fmt.Sprintf(msg, a...)
	_, fileName, fileLine, ok := runtime.Caller(2)
	if (ok) {
		log.Printf("%s at %s:%d\n", outmsg, fileName, fileLine)
	} else {
		log.Printf("%s\n", outmsg)
	}
	t.FailNow()
}

func isok(t *testing.T, err error) {
	if err != nil {
		failCaller2(t, "Error! Error is not nil! - %s", err)
	}
}

func iseof(t *testing.T, err error) {
	if err != io.EOF {
		failCaller2(t, "Error! Error is not EOF! - %s", err)
	}
}

func assert(t *testing.T, expr bool) {
	if !expr {
		failCaller2(t, "Assertion failed!")
	}
}

func testFoobarAr(reader io.Reader, t *testing.T) {
	ar, err := deb.LoadAr(reader)
	isok(t, err)

	foo, err := ar.Next()
	isok(t, err)
	assert(t, foo.Name      == "foo.txt")
	assert(t, foo.Timestamp == 0)
	assert(t, foo.OwnerID   == 0)
	assert(t, foo.GroupID   == 0)
	assert(t, foo.FileMode  == "644")
	assert(t, foo.Size      == 7)
	part := make([]byte, 4)
	n, err := foo.Data.Read(part)
	isok(t, err)
	assert(t, n <= 4)
	assert(t, string(part[:n]) == "FooB"[:n])

	bar, err := ar.Next()
	isok(t, err)
	assert(t, bar.Name      == "bar.txt")
	assert(t, bar.Timestamp == 1559837289)
	assert(t, bar.OwnerID   == 1000)
	assert(t, bar.GroupID   == 1000)
	assert(t, bar.FileMode  == "100644")
	assert(t, bar.Size      == 21)
	data, err := ioutil.ReadAll(bar.Data)
	isok(t, err)
	assert(t, string(data) == "FooBar\nFooBar\nFooBar\n")

	na, err := ar.Next()
	iseof(t, err)
	assert(t, na == nil)
}

func TestFoobarAr(t *testing.T) {
	foobar_b := []byte(foobar_ar)
	testFoobarAr(&wholeReader{Content: foobar_b}, t)
	testFoobarAr(&byteReader{Content: foobar_b}, t)
}
