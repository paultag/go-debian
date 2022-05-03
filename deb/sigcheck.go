package deb // import "pault.ag/go/debian/deb"

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/openpgp"
)

const (
	SigTypeArchive = `archive`
	SigTypeMaint   = `maint`
	SigTypeOrigin  = `origin`
)

func (deb *Deb) CheckDebsig(validKeys openpgp.EntityList, sigType string) (signer *openpgp.Entity, err error) {
	sig, ok := deb.ArContent[`_gpg`+sigType]
	if !ok {
		return nil, fmt.Errorf("no signature of type %v present", sigType)
	}

	binaryFlag, ok := deb.ArContent[`debian-binary`]
	if !ok {
		return nil, fmt.Errorf("archive does not contain a debian-binary flag")
	}

	var control, data *ArEntry
	for _, member := range deb.ArContent {
		if strings.HasPrefix(member.Name, "control.") {
			control = member
			if data != nil {
				break
			}
		} else if strings.HasPrefix(member.Name, "data.") {
			data = member
			if control != nil {
				break
			}
		}
	}
	if control == nil || data == nil {
		return nil, fmt.Errorf("unable to find signed data")
	}
	binaryFlag.Data.Seek(0, 0)
	control.Data.Seek(0, 0)
	data.Data.Seek(0, 0)
	signedData := io.MultiReader(binaryFlag.Data, control.Data, data.Data)
	return openpgp.CheckDetachedSignature(validKeys, signedData, sig.Data)
}
