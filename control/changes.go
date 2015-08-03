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
	"pault.ag/go/debian/version"
)

type ChangesFileHash struct {
	// cb136f28a8c971d4299cc68e8fdad93a8ca7daf3 1131 dput-ng_1.9.dsc
	Algorithm string
	Hash      string
	Size      int
	Filename  string
}

type Changes struct {
	Paragraph

	Filename string

	Format          string
	Source          string
	Binaries        []string
	Architectures   []dependency.Arch
	Version         version.Version
	Origin          string
	Distribution    string
	Urgency         string
	Maintainer      string
	ChangedBy       string
	Closes          []string
	Changes         string
	ChecksumsSha1   []ChangesFileHash
	ChecksumsSha256 []ChangesFileHash
	Files           []ChangesFileHash
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

func ParseChanges(reader *bufio.Reader, path string) (ret *Changes, err error) {

	/* a Changes is a Paragraph, with some stuff. So, let's first take
	 * the bufio.Reader and produce a stock Paragraph. */
	src, err := ParseParagraph(reader)
	if err != nil {
		return nil, err
	}

	arch, err := dependency.ParseArchitectures(src.Values["Architecture"])
	if err != nil {
		return nil, err
	}

	version, err := version.Parse(src.Values["Version"])
	if err != nil {
		return nil, err
	}

	ret = &Changes{
		Paragraph: *src,
		Filename:  path,

		Format:     src.Values["Format"],
		Source:     src.Values["Source"],
		Origin:     src.Values["Origin"],
		Urgency:    src.Values["Urgency"],
		Maintainer: src.Values["Maintainer"],
		ChangedBy:  src.Values["Changed-By"],
		Changes:    src.Values["Changes"],

		Architectures: arch,
		Version:       version,
		Distribution:  src.Values["Distribution"],

		ChecksumsSha1:   parseHashes(src.Values["Checksums-Sha1"], "SHA1"),
		ChecksumsSha256: parseHashes(src.Values["Checksums-Sha256"], "SHA256"),
		Files:           parseHashes(src.Values["Files"], "MD5"),
		Binaries:        strings.Split(src.Values["Binary"], " "),
		Closes:          strings.Split(src.Values["Closes"], " "),
	}

	return
}

func (changes *Changes) GetDSC() (*DSC, error) {
	for _, file := range changes.Files {
		if strings.HasSuffix(file.Filename, ".dsc") {
			dsc, err := ParseDscFile(file.Filename)
			if err != nil {
				return nil, err
			}
			return dsc, nil
		}
	}
	return nil, fmt.Errorf("No .dsc file in .changes")
}

func (changes *Changes) Move(dest string) error {
	if file, err := os.Stat(dest); err == nil && !file.IsDir() {
		return fmt.Errorf("Attempting to move .changes to a non-directory")
	}
	dirname := filepath.Base(changes.Filename)

	changes.Filename = dest + "/" + dirname
	for _, file := range changes.Files {
		dirname := filepath.Base(changes.Filename)
		err := os.Rename(file.Filename, dest+"/"+dirname)
		if err != nil {
			return err
		}
	}

	err := os.Rename(changes.Filename, dest+"/"+dirname)
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
