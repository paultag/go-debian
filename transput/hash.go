package transput

import (
	"fmt"
	"hash"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

func GetHash(name string) (hash.Hash, error) {
	switch name {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("Unknown algorithm: %s", name)
	}
}

func NewHashWriter(name string) (*HashWriter, error) {
	hash, err := GetHash(name)
	if err != nil {
		return nil, err
	}

	hw := HashWriter{
		hash: hash,
		size: 0,
	}

	return &hw, nil
}

type HashWriter struct {
	hash hash.Hash
	size int64
}

func (dh *HashWriter) Write(p []byte) (int, error) {
	n, err := dh.hash.Write(p)
	dh.size += int64(n)
	return n, err
}

func (dh *HashWriter) Size() int64 {
	return dh.size
}

func (dh *HashWriter) Sum(b []byte) []byte {
	return dh.hash.Sum(b)
}
