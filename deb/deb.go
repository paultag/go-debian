package deb

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DebEntry struct {
	Name      string
	Timestamp int64
	OwnerID   int64
	GroupID   int64
	FileMode  string
	Size      int64
	Data      io.Reader
}

type Deb struct {
	file       *os.File
	jumpTarget int64
}

func (d *Deb) jumpAround() error {
	_, err := d.file.Seek(d.jumpTarget, 0)
	return err
}

func (d *Deb) setJumpTarget(e *DebEntry) error {
	place, err := d.file.Seek(0, 1)
	target := place + e.Size
	if target%1 == 1 {
		target += 1
	}
	d.jumpTarget = target
	return err
}

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

func toDecimal(input []byte) (int64, error) {
	stream := strings.TrimSpace(string(input))
	out, err := strconv.Atoi(stream)
	return int64(out), err
}

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

func Load(path string) (*Deb, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return createDeb(fd)
}

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
