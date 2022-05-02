package deb_test

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/pschou/go_debian/deb"
)

/*
 *
 */

func isok(t *testing.T, err error) {
	if err != nil && err != io.EOF {
		log.Printf("Error! Error is not nil! %s\n", err)
		t.FailNow()
	}
}

func notok(t *testing.T, err error) {
	if err == nil {
		log.Printf("Error! Error is nil!\n")
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

// As discussed in testdata/README.md, `multi_archive.a` is taken from Blake
// Smith's ar project. Some of the data below follows.

func TestAr(t *testing.T) {
	file, err := os.Open("testdata/multi_archive.a")
	isok(t, err)

	ar, err := deb.LoadAr(file)
	isok(t, err)

	firstEntry, err := ar.Next()
	isok(t, err)

	assert(t, firstEntry.Name == `hello.txt`)
	assert(t, firstEntry.Timestamp == 1361157466)
	assert(t, firstEntry.OwnerID == 501)
	assert(t, firstEntry.GroupID == 20)

	firstContent, err := ioutil.ReadAll(firstEntry.Data)
	isok(t, err)
	assert(t, firstEntry.Size == int64(len(firstContent)))
	assert(t, string(firstContent) == "Hello world!\n")

	secondEntry, err := ar.Next()
	isok(t, err)
	secondContent, err := ioutil.ReadAll(secondEntry.Data)
	isok(t, err)
	assert(t, secondEntry.Size == int64(len(secondContent)))
	assert(t, string(secondContent) == "I love lamp.\n")

	// Now, test that we can rewind and reread the first file even after
	// reading the second one.
	_, err = firstEntry.Data.Seek(0, 0)
	isok(t, err)
	firstRereadContent, err := ioutil.ReadAll(firstEntry.Data)
	isok(t, err)
	assert(t, string(firstContent) == string(firstRereadContent))
}
