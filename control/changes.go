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

package control

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pault.ag/go/debian/dependency"
	"pault.ag/go/debian/internal"
	"pault.ag/go/debian/version"
)

type ChangesFileHash struct {
	// cb136f28a8c971d4299cc68e8fdad93a8ca7daf3 1131 dput-ng_1.9.dsc
	Algorithm string
	Hash      string
	Size      int
	Filename  string
}

type FileListChangesFileHash struct {
	ChangesFileHash

	Component string
	Priority  string
}

func (c *FileListChangesFileHash) UnmarshalControl(data string) error {
	var err error
	c.Algorithm = "md5"
	vals := strings.Split(data, " ")
	if len(data) < 5 {
		return fmt.Errorf("Error: Unknown File List Hash line: '%s'", data)
	}

	c.Hash = vals[0]
	c.Size, err = strconv.Atoi(vals[1])
	if err != nil {
		return err
	}
	c.Component = vals[2]
	c.Priority = vals[3]

	c.Filename = vals[4]
	return nil
}

type SHAChangesFileHash struct {
	ChangesFileHash
}

func (c *SHAChangesFileHash) unmarshalControl(algorithm, data string) error {
	var err error
	c.Algorithm = algorithm
	vals := strings.Split(data, " ")
	if len(data) < 4 {
		return fmt.Errorf("Error: Unknown SHA Hash line: '%s'", data)
	}

	c.Hash = vals[0]
	c.Size, err = strconv.Atoi(vals[1])
	if err != nil {
		return err
	}
	c.Filename = vals[2]
	return nil
}

type SHA1ChangesFileHash struct{ SHAChangesFileHash }

func (c *SHA1ChangesFileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("SHA1", data)
}

type SHA256ChangesFileHash struct{ SHAChangesFileHash }

func (c *SHA256ChangesFileHash) UnmarshalControl(data string) error {
	return c.unmarshalControl("SHA256", data)
}

type Changes struct {
	Paragraph

	Filename string

	Format          string
	Source          string
	Binaries        []string          `control:"Binary" delim:" "`
	Architectures   []dependency.Arch `control:"Architecture"`
	Version         version.Version
	Origin          string
	Distribution    string
	Urgency         string
	Maintainer      string
	ChangedBy       string `control:"Changed-By"`
	Closes          []string
	Changes         string
	ChecksumsSha1   []SHA1ChangesFileHash     `control:"Checksums-Sha1" delim:"\n" strip:"\n\r\t "`
	ChecksumsSha256 []SHA256ChangesFileHash   `control:"Checksums-Sha256" delim:"\n" strip:"\n\r\t "`
	Files           []FileListChangesFileHash `control:"Files" delim:"\n" strip:"\n\r\t "`
}

func parseHashes(buf string, algorithm string) (ret []ChangesFileHash) {
	for _, el := range strings.Split(buf, "\n") {
		if el == "" {
			continue
		}
		vals := strings.Split(el, " ")

		/* Sanity check length here */
		size, err := strconv.Atoi(vals[1])
		if err != nil {
			continue
		}

		switch len(vals) {
		case 5:
			ret = append(ret, ChangesFileHash{
				Algorithm: algorithm,
				Hash:      vals[0],
				Size:      size,
				Filename:  vals[4],
			})
		case 3:
			ret = append(ret, ChangesFileHash{
				Algorithm: algorithm,
				Hash:      vals[0],
				Size:      size,
				Filename:  vals[2],
			})
		}
	}
	return
}

func ParseChangesFile(path string) (ret *Changes, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ret, err = ParseChanges(bufio.NewReader(f), path)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ParseChanges(reader *bufio.Reader, path string) (*Changes, error) {
	ret := &Changes{Filename: path}
	if err := Unmarshal(ret, reader); err != nil {
		return nil, err
	}
	return ret, nil
}

func (changes *Changes) GetDSC() (*DSC, error) {
	for _, file := range changes.Files {
		if strings.HasSuffix(file.Filename, ".dsc") {

			// Right, now lets resolve the absolute path.
			baseDir := filepath.Dir(changes.Filename)

			dsc, err := ParseDscFile(baseDir + "/" + file.Filename)
			if err != nil {
				return nil, err
			}
			return dsc, nil
		}
	}
	return nil, fmt.Errorf("No .dsc file in .changes")
}

func (changes *Changes) Copy(dest string) error {
	if file, err := os.Stat(dest); err == nil && !file.IsDir() {
		return fmt.Errorf("Attempting to move .changes to a non-directory")
	}

	for _, file := range changes.Files {
		dirname := filepath.Base(file.Filename)
		err := internal.Copy(file.Filename, dest+"/"+dirname)
		if err != nil {
			return err
		}
	}

	dirname := filepath.Base(changes.Filename)
	err := internal.Copy(changes.Filename, dest+"/"+dirname)
	changes.Filename = dest + "/" + dirname

	if err != nil {
		return err
	}

	return nil
}

func (changes *Changes) Move(dest string) error {
	if file, err := os.Stat(dest); err == nil && !file.IsDir() {
		return fmt.Errorf("Attempting to move .changes to a non-directory")
	}

	for _, file := range changes.Files {
		dirname := filepath.Base(file.Filename)
		err := os.Rename(file.Filename, dest+"/"+dirname)
		if err != nil {
			return err
		}
	}

	dirname := filepath.Base(changes.Filename)
	err := os.Rename(changes.Filename, dest+"/"+dirname)
	changes.Filename = dest + "/" + dirname

	if err != nil {
		return err
	}

	return nil
}

func (changes *Changes) Remove() error {
	for _, file := range changes.Files {
		err := os.Remove(file.Filename)
		if err != nil {
			return err
		}
	}
	err := os.Remove(changes.Filename)
	if err != nil {
		return err
	}
	return nil
}

// vim: foldmethod=marker
