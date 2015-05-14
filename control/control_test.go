package control_test

import (
	"bufio"
	"strings"
	"testing"

	"pault.ag/x/go-debian/control"
)

/*
 *
 */

func TestDependencyControlParse(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader(`Source: fbautostart
Section: misc
Priority: optional
Maintainer: Paul Tagliamonte <paultag@ubuntu.com>
Build-Depends: debhelper (>= 9)
Standards-Version: 3.9.3
Homepage: https://launchpad.net/fbautostart
Vcs-Git: git://git.debian.org/collab-maint/fbautostart.git
Vcs-Browser: http://git.debian.org/?p=collab-maint/fbautostart.git

Package: fbautostart
Architecture: any
Depends: ${shlibs:Depends}, ${misc:Depends}
Description: XDG compliant autostarting app for Fluxbox
 The fbautostart app was designed to have little to no overhead, while
 still maintaining the needed functionality of launching applications
 according to the XDG spec.
 .
 This package contains support for GNOME and KDE.
`))
	c, err := control.ParseControl(reader)
	isok(t, err)
	assert(t, c != nil)
	assert(t, len(c.Binaries) == 1)

	depends, err := c.Source.GetBuildDepends()
	isok(t, err)

	assert(t, depends.Relations[0].Possibilities[0].Name == "debhelper")
	assert(t, depends.Relations[0].Possibilities[0].Version.Number == "9")
	assert(t, depends.Relations[0].Possibilities[0].Version.Operator == ">=")

	assert(t, c.Binaries[0].Arch.CPU == "any")
}
