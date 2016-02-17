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

// ArEntry {{{

// Container type to access the different parts of a Debian `ar(1)` Archive.
//
// The most interesting parts of this are the `Name` attribute, Data
// `io.Reader`, and the Tarfile helpers. This will allow the developer to
// programatically inspect the information inside without forcing her to
// unpack the .deb to the filesystem.
type ArEntry struct {
	Name      string
	Timestamp int64
	OwnerID   int64
	GroupID   int64
	FileMode  string
	Size      int64
	Data      io.Reader
}

// }}}

// Ar {{{

// This struct encapsulates a Debian .deb flavored `ar(1)` archive.
type Ar struct {
	file       *os.File
	jumpTarget int64
}

// Close the underlying os.File object.
func (d *Ar) Close() error {
	return d.file.Close()
}

// Given a path on the filesystem (`path`), return a `Ar` struct to allow
// for programatic access to its contents.
func LoadAr(path string) (*Ar, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return createAr(fd)
}

// Helpers {{{

// This will reset the underlying File offset back to the next Entry,
// without having to worry about how much data the user read off
// the Data attribute.
func (d *Ar) jumpAround() error {
	_, err := d.file.Seek(d.jumpTarget, 0)
	return err
}

// This will calculate the next jump target to be used after we regain
// control, and need to give the user the next entry.
func (d *Ar) setJumpTarget(e *ArEntry) error {
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
func (d *Ar) Next() (*ArEntry, error) {
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
	entry, err := parseArEntry(line)
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

// parseArEntry {{{

// Take the AR format line, and create a ArEntry (without .Data set)
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
func parseArEntry(line []byte) (*ArEntry, error) {
	if len(line) != 60 {
		return nil, fmt.Errorf("Malformed file entry line length")
	}

	if line[58] != 0x60 && line[59] != 0x0A {
		return nil, fmt.Errorf("Malformed file entry line endings")
	}

	entry := ArEntry{
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
		return fmt.Errorf("Header does't look right!")
	}
	return nil
}

// }}}

// createAr {{{

// Given a brand spank'n new os.File entry, go ahead and check that it
// looks authentic, and create an Ar struct. This is basically what `Load`
// does.
func createAr(reader *os.File) (*Ar, error) {
	if err := checkAr(reader); err != nil {
		return nil, err
	}
	debFile := Ar{
		file:       reader,
		jumpTarget: 8, /* See the !<arch>\n header above */
	}
	return &debFile, nil
}

// }}}

// }}}

// vim: foldmethod=marker
