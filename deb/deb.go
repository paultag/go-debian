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

package deb

import (
	"fmt"
	"path"
	"strings"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/version"
)

// Control {{{

// Binary Control format, as exists in the Control section of the `.deb`
// archive, as defined in Debian Policy, section 5.3, entitiled "Binary
// package control files -- DEBIAN/control".
type Control struct {
	control.Paragraph

	Package       string          `required:"true"`
	Version       version.Version `required:"true"`
	Architecture  dependency.Arch `required:"true"`
	Maintainer    string          `required:"true"`
	InstalledSize int             `control:"Installed-Size"`
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

// }}}

// Deb {{{

// Container struct to encapsulate a `.deb` file on disk. This contains
// information about what exactly we're looking at. When loaded. information
// regarding the Control file is read from the control section of the .deb,
// and Unmarshaled into the `Control` member of the Struct.
type Deb struct {
	Control Control
	Path    string
}

// Load {{{

// Load a given `.deb` off disk and into a `Deb` container struct.
func Load(pathname string) (*Deb, error) {
	ar, err := LoadAr(pathname)
	if err != nil {
		return nil, err
	}

	defer ar.Close()

	var controlEntry *ArEntry

	for {
		entry, err := ar.Next()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(entry.Name, "control.") && entry.IsTarfile() {
			controlEntry = entry
			break
		}
	}

	if controlEntry == nil {
		return nil, fmt.Errorf("No control blob found!")
	}

	tarFile, err := controlEntry.Tarfile()
	if err != nil {
		return nil, err
	}

	/* Now, scan for control */
	for {
		tfEntry, err := tarFile.Next()
		if err != nil {
			return nil, err
		}
		if path.Clean(tfEntry.Name) == "control" {
			break
		}
	}

	var debControl = Control{}
	if err := control.Unmarshal(&debControl, tarFile); err != nil {
		return nil, err
	}
	deb := Deb{
		Control: debControl,
		Path:    pathname,
	}
	return &deb, nil
}

// }}}

// }}}

// vim: foldmethod=marker
