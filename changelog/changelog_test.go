package changelog_test

import (
	"bufio"
	"log"
	"strings"
	"testing"

	"pault.ag/go/debian/changelog"
)

/*
 *
 */

func isok(t *testing.T, err error) {
	if err != nil {
		log.Printf("Error! Error is not nil! %s\n", err)
		t.FailNow()
	}
}

func notok(t *testing.T, err error) {
	if err == nil {
		log.Printf("Error! Error is nil!\n")
		t.FailNow()
	}
}

func assert(t *testing.T, expr bool) {
	if !expr {
		log.Printf("Assertion failed!")
		t.FailNow()
	}
}

/*
 *
 */

// {{{ test changelog entry
var changeLog = `hello (2.10-1) unstable; urgency=low

  * New upstream release.
  * debian/patches: Drop 01-fix-i18n-of-default-message, no longer needed.
  * debian/patches: Drop 99-config-guess-config-sub, no longer needed.
  * debian/rules: Drop override_dh_auto_build hack, no longer needed.
  * Standards-Version: 3.9.6 (no changes for this).

 -- Santiago Vila <sanvila@debian.org>  Sun, 22 Mar 2015 11:56:00 +0100

hello (2.9-2) unstable; urgency=low

  * Apply patch from Reuben Thomas to fix i18n of default message.
    This is upstream commit c4aed00. Closes: #767172.
  * The previous change in src/hello.c trigger a rebuild of man/hello.1
    that we don't need. Add a "touch man/hello.1" to avoid it.
  * Use Breaks: hello-debhelper (<< 2.9), not Conflicts,
    as hello-debhelper is deprecated.
  * Restore simple watch file from old hello package that was lost
    when the packages were renamed.
  * Update 99-config-guess-config-sub patch.

 -- Santiago Vila <sanvila@debian.org>  Thu, 06 Nov 2014 12:03:40 +0100
`

// }}}

func TestChangelogEntry(t *testing.T) {
	changeLog, err := changelog.ParseOne(bufio.NewReader(strings.NewReader(changeLog)))
	isok(t, err)
	assert(t, changeLog.ChangedBy == "Santiago Vila <sanvila@debian.org>")
}

func TestChangelogEntries(t *testing.T) {
	changeLogs, err := changelog.Parse(strings.NewReader(changeLog))
	isok(t, err)
	assert(t, len(changeLogs) == 2)
}

// vim: foldmethod=marker
