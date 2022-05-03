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

package deb // import "pault.ag/go/debian/deb"

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

// Control {{{

// Binary Control format, as exists in the Control section of the `.deb`
// archive, as defined in Debian Policy, section 5.3, entitled "Binary
// package control files -- DEBIAN/control".
type Control struct {
	control.Paragraph

	Package       string `required:"true"`
	Source        string
	Version       version.Version `required:"true"`
	Architecture  dependency.Arch `required:"true"`
	Maintainer    string          `required:"true"`
	InstalledSize int             `control:"Installed-Size"`
	MultiArch     string          `control:"Multi-Arch"`
	Depends       dependency.Dependency
	Recommends    dependency.Dependency
	Suggests      dependency.Dependency
	Breaks        dependency.Dependency
	Replaces      dependency.Dependency
	BuiltUsing    dependency.Dependency `control:"Built-Using"`
	Section       string
	Priority      string
	Homepage      string
	Description   string `required:"true"`
}

func (c Control) SourceName() string {
	if c.Source == "" {
		return c.Package
	}
	return c.Source
}

// }}}

// Deb {{{

// Container struct to encapsulate a `.deb` file on disk. This contains
// information about what exactly we're looking at. When loaded. information
// regarding the Control file is read from the control section of the .deb,
// and Unmarshaled into the `Control` member of the Struct.
type Deb struct {
	Control    Control
	Path       string
	Data       *tar.Reader
	Closer     io.Closer
	ControlExt string
	DataExt    string
	ArContent  map[string]*ArEntry
}

func (deb *Deb) Close() error {
	if deb.Closer != nil {
		return deb.Closer.Close()
	}
	return nil
}

// Load {{{

// Load {{{

// Given a reader, and the file path to the file (for use in the Deb later)
// create a deb.Deb object, and populate the Control and Data members.
// It is the caller's responsibility to call Close() when done.
func Load(in io.ReaderAt, pathname string) (*Deb, error) {
	ar, err := LoadAr(in)
	if err != nil {
		return nil, err
	}
	deb, err := loadDeb(ar)
	if err != nil {
		return nil, err
	}
	deb.Path = pathname
	return deb, nil
}

// }}}

// LoadFile {{{

type Closer func() error

type closerAdapter struct {
	closeFunc Closer
}

func (c *closerAdapter) Close() error {
	return c.closeFunc()
}

func LoadFile(path string) (*Deb, Closer, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	debFile, err := Load(fd, path)
	if err != nil {
		fd.Close()
		return nil, nil, err
	}

	dataCloser := debFile.Closer

	// Replace debFile.Closer function with another one that also closes fd.
	// We do this to preserve backwards compatibility, ensuring users already invoking the returned
	// closeFunc do not need to also call Close() on debFile. Earlier versions of this library did
	// not require Close() to be invoked on debFile.
	closeFunc := func() error {
		err1 := dataCloser.Close()
		err2 := fd.Close()
		if err1 != nil {
			return err1
		}
		return err2
	}

	debFile.Closer = &closerAdapter{closeFunc}

	return debFile, closeFunc, nil

}

// }}}

// Debian .deb Loader Internals {{{

// Top-level .deb loader dispatch on Version {{{

// Look for the debian-binary member and figure out which version to read
// it as. Return the newly created .deb struct.
func loadDeb(archive *Ar) (*Deb, error) {
	contents := make(map[string]*ArEntry)
	for {
		member, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		contents[member.Name] = member
	}
	member, ok := contents["debian-binary"]
	if !ok {
		return nil, fmt.Errorf("Archive contains no binary version member!")
	}
	reader := bufio.NewReader(member.Data)
	version, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	switch version {
	case "2.0\n":
		return loadDeb2(contents)
	default:
		return nil, fmt.Errorf("Unknown binary version: '%s'", version)
	}
}

// }}}

// Debian .deb format 2.0 {{{

// Top-level .deb loader dispatch for 2.0 {{{

// Load a Debian 2.x series .deb - track down the control and data members.
func loadDeb2(archive map[string]*ArEntry) (*Deb, error) {
	ret := Deb{ArContent: archive}

	if err := loadDeb2Control(archive, &ret); err != nil {
		return nil, err
	}

	if err := loadDeb2Data(archive, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// }}}

// Decode .deb 2.0 control data into the struct {{{

// Load a Debian 2.x series .deb control file and write it out to
// the deb.Deb.Control member.
func loadDeb2Control(archive map[string]*ArEntry, deb *Deb) error {
	for _, member := range archive {
		if strings.HasPrefix(member.Name, "control.") {
			archive, closer, err := member.Tarfile()
			if err != nil {
				return err
			}
			deb.ControlExt = member.Name[8:len(member.Name)]
			for {
				member, err := archive.Next()
				if err != nil {
					closer.Close()
					return err
				}
				if path.Clean(member.Name) == "control" {
					err1 := control.Unmarshal(&deb.Control, archive)
					err2 := closer.Close()
					if err1 != nil {
						return err1
					}
					return err2
				}
			}
			closer.Close()
		}
	}
	return fmt.Errorf("Missing or out of order .deb member 'control'")
}

// }}}

// Decode .deb 2.0 package data into the struct {{{

// Load a Debian 2.x series .deb data file and write it out to
// the deb.Deb.Data member.
func loadDeb2Data(archive map[string]*ArEntry, deb *Deb) error {
	for _, member := range archive {
		if strings.HasPrefix(member.Name, "data.") {
			archive, closer, err := member.Tarfile()
			if err != nil {
				return err
			}
			deb.DataExt = member.Name[5:len(member.Name)]
			deb.Data = archive
			deb.Closer = closer
			return nil
		}
	}
	return fmt.Errorf("Missing or out of order .deb member 'data'")
}

// }}}

// }}} }}} }}} }}}

// vim: foldmethod=marker
