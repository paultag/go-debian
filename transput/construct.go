package transput

import (
	"io"
)

func NewHasherWriter(hash string, target io.Writer) (io.Writer, *Hasher, error) {
	hw, err := NewHasher(hash)
	if err != nil {
		return nil, nil, err
	}
	endWriter := io.MultiWriter(target, hw)
	return endWriter, hw, nil
}

func NewHasherReader(hash string, target io.Reader) (io.Reader, *Hasher, error) {
	hw, err := NewHasher(hash)
	if err != nil {
		return nil, nil, err
	}
	endReader := io.TeeReader(target, hw)
	return endReader, hw, nil
}

func NewCompressedHasherWriter(hash, compression string, target io.Writer) (io.WriteCloser, *Hasher, error) {
	w, h, err := NewHasherWriter(hash, target)
	if err != nil {
		return nil, nil, err
	}
	compressor, err := GetCompressor(compression)
	if err != nil {
		return nil, nil, err
	}
	newWriter, err := compressor(w)
	if err != nil {
		return nil, nil, err
	}
	return newWriter, h, nil
}
