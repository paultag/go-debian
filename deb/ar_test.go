package deb_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"pault.ag/go/debian/deb"
)

/*
 *
 */

func isok(t *testing.T, err error) {
	t.Helper()
	if err != nil && err != io.EOF {
		t.Fatalf("Error! Error is not nil! %v", err)
	}
}

func notok(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Error! Error is nil!")
	}
}

func assert(t *testing.T, expr bool) {
	t.Helper()
	if !expr {
		t.Fatal("Assertion failed!")
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

	firstContent, err := io.ReadAll(firstEntry.Data)
	isok(t, err)
	assert(t, firstEntry.Size == int64(len(firstContent)))
	assert(t, string(firstContent) == "Hello world!\n")

	secondEntry, err := ar.Next()
	isok(t, err)
	secondContent, err := io.ReadAll(secondEntry.Data)
	isok(t, err)
	assert(t, secondEntry.Size == int64(len(secondContent)))
	assert(t, string(secondContent) == "I love lamp.\n")

	// Now, test that we can rewind and reread the first file even after
	// reading the second one.
	_, err = firstEntry.Data.Seek(0, 0)
	isok(t, err)
	firstRereadContent, err := io.ReadAll(firstEntry.Data)
	isok(t, err)
	assert(t, string(firstContent) == string(firstRereadContent))
}

// `long.a` is taken from gccgoexportdata.
// It contains empty fields.
func TestArEmptyFields(t *testing.T) {
	file, err := os.Open("testdata/long.a")
	isok(t, err)

	ar, err := deb.LoadAr(file)
	isok(t, err)

	// First
	firstEntry, err := ar.Next()
	isok(t, err)

	assert(t, firstEntry.Name == ``)
	assert(t, firstEntry.Timestamp == 1478713642)
	assert(t, firstEntry.OwnerID == 0)
	assert(t, firstEntry.GroupID == 0)
	assert(t, firstEntry.Size == 4)

	firstContent, err := io.ReadAll(firstEntry.Data)
	isok(t, err)
	assert(t, firstEntry.Size == int64(len(firstContent)))
	//assert(t, string(firstContent) == "")

	// Second
	secondEntry, err := ar.Next()
	isok(t, err)

	assert(t, secondEntry.Name == `/`)
	assert(t, secondEntry.Timestamp == 0)
	assert(t, secondEntry.OwnerID == 0)
	assert(t, secondEntry.GroupID == 0)
	assert(t, secondEntry.Size == 32)

	secondContent, err := io.ReadAll(secondEntry.Data)
	isok(t, err)
	assert(t, secondEntry.Size == int64(len(secondContent)))
	assert(t, string(secondContent) == "name-longer-than-16-bytes.gox/\n\n")

	// Third
	thirdEntry, err := ar.Next()
	isok(t, err)
	fmt.Printf("%v", thirdEntry)
	//// Now, test that we can rewind and reread the first file even after
	//// reading the second one.
	//_, err = firstEntry.Data.Seek(0, 0)
	//isok(t, err)
	//firstRereadContent, err := io.ReadAll(firstEntry.Data)
	//isok(t, err)
	//assert(t, string(firstContent) == string(firstRereadContent), "")
}
