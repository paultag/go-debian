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
	"io"
	"os"
	"strconv"
	"strings"
)

// DebEntry {{{

// Container type to access the different parts of a Debian binary package.
//
// The most interesting parts of this are the `Name` attribute, Data
// `io.Reader`, and the Tarfile helpers. This will allow the developer to
// programatically inspect the information inside without forcing her to
// unpack the .deb to the filesystem.
type DebEntry struct {
	Name      string
	Timestamp int64
	OwnerID   int64
	GroupID   int64
	FileMode  string
	Size      int64
	Data      io.Reader
}

// }}}

// Deb {{{

// A Debian .deb file is, at a high level, an `ar(1)` archive, which contains
// (typically 3) members. A file to note which format is being used
// (`debian-binary`), related Debian control files `control.tar.*`, and the
// actual contents of the Binary distribution `data.tar.*`.
type Deb struct {
	file       *os.File
	jumpTarget int64
}

func Load(path string) (*Deb, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return createDeb(fd)
}

// Helpers {{{

// This will reset the underlying File offset back to the next Entry,
// without having to worry about how much data the user read off
// the Data attribute.
func (d *Deb) jumpAround() error {
	_, err := d.file.Seek(d.jumpTarget, 0)
	return err
}

// This will calculate the next jump target to be used after we regain
// control, and need to give the user the next entry.
func (d *Deb) setJumpTarget(e *DebEntry) error {
	place, err := d.file.Seek(0, 1)
	target := place + e.Size
	if target%1 == 1 {
		target += 1
	}
	d.jumpTarget = target
	return err
}

// }}}

// Next {{{

// Function to jump to the next file in the Debian `ar(1)` archive, and
// return the next member.
func (d *Deb) Next() (*DebEntry, error) {
	line := make([]byte, 60)

	if err := d.jumpAround(); err != nil {
		return nil, err
	}

	count, err := d.file.Read(line)
	if err != nil {
		return nil, err
	}
	if count == 1 && line[0] == '\n' {
		return nil, io.EOF
	}
	if count != 60 {
		return nil, fmt.Errorf("Caught a short read at the end")
	}
	entry, err := parseDebEntry(line)
	if err != nil {
		return nil, err
	}

	if err := d.setJumpTarget(entry); err != nil {
		return nil, err
	}

	entry.Data = io.LimitReader(d.file, entry.Size)

	return entry, nil
}

// }}}

// toDecimal {{{

// Take a byte array, and return an int64
func toDecimal(input []byte) (int64, error) {
	stream := strings.TrimSpace(string(input))
	out, err := strconv.Atoi(stream)
	return int64(out), err
}

// }}}

// }}}

// AR Format Hackery {{{

// parseDebEntry {{{

// Take the AR format line, and create a DebEntry (without .Data set)
// to be returned to the user later.
//
// +-------------------------------------------------------
// | Offset  Length  Name                         Format
// +-------------------------------------------------------
// | 0       16      File name                    ASCII
// | 16      12      File modification timestamp  Decimal
// | 28      6       Owner ID                     Decimal
// | 34      6       Group ID                     Decimal
// | 40      8       File mode                    Octal
// | 48      10      File size in bytes           Decimal
// | 58      2       File magic                   0x60 0x0A
//
func parseDebEntry(line []byte) (*DebEntry, error) {
	if len(line) != 60 {
		return nil, fmt.Errorf("Malformed file entry line length")
	}

	if line[58] != 0x60 && line[59] != 0x0A {
		return nil, fmt.Errorf("Malformed file entry line endings")
	}

	entry := DebEntry{
		Name:     strings.TrimSpace(string(line[0:16])),
		FileMode: strings.TrimSpace(string(line[48:58])),
	}

	for target, value := range map[*int64][]byte{
		&entry.Timestamp: line[16:28],
		&entry.OwnerID:   line[28:34],
		&entry.GroupID:   line[34:40],
		&entry.Size:      line[48:58],
	} {
		intValue, err := toDecimal(value)
		if err != nil {
			return nil, err
		}
		*target = intValue
	}

	return &entry, nil
}

// }}}

// checkAr {{{

// Given a brand spank'n new os.File entry, go ahead and make sure it looks
// like an `ar(1)` archive, and not some rando file.
func checkAr(reader *os.File) error {
	header := make([]byte, 8)
	if _, err := reader.Read(header); err != nil {
		return err
	}
	if string(header) != "!<arch>\n" {
		return fmt.Errorf("Header doens't look right!")
	}
	return nil
}

// }}}

// createDeb {{{

// Given a brand spank'n new os.File entry, go ahead and check that it
// looks authentic, and create a Deb. This is basically what `Load` does.
func createDeb(reader *os.File) (*Deb, error) {
	if err := checkAr(reader); err != nil {
		return nil, err
	}
	debFile := Deb{
		file:       reader,
		jumpTarget: 8, /* See the !<arch>\n header above */
	}
	return &debFile, nil
}

// }}}

// }}}

// vim: foldmethod=marker
