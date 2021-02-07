package deb_test

import (
	"golang.org/x/crypto/openpgp"
	"os"
	"testing"

	"pault.ag/go/debian/deb"
)

func TestCheckDebsig(t *testing.T) {
	testDeb2, closer, err := deb.LoadFile(`testdata/sigtest1_1.0-1_all.deb`)
	isok(t, err)
	assert(t, closer != nil)
	defer closer()
	keyring, err := os.Open("testdata/keyrings/7CD73F641E04EC2D/debian-debsig.gpg")
	isok(t, err)
	keys, err := openpgp.ReadKeyRing(keyring)
	isok(t, err)
	entityList, err := testDeb2.CheckDebsig(keys, deb.SigTypeOrigin)
	isok(t, err)
	assert(t, entityList != nil)
}
