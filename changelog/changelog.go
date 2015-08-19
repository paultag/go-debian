package changelog

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ChangelogEntry struct {
	Source    string
	Version   string
	Target    string
	Arguments map[string]string
	Changelog string
	ChangedBy string
	When      string
}

func trim(line string) string {
	return strings.Trim(line, "\n\r\t ")
}

func partition(line, delim string) (string, string) {
	entries := strings.SplitN(line, delim, 2)
	if len(entries) != 2 {
		return line, ""
	}
	return entries[0], entries[1]

}

func ParseOne(reader *bufio.Reader) (*ChangelogEntry, error) {
	changeLog := ChangelogEntry{}

	var header string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\n" {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			/* Great. Let's work with this. */
			header = line
			break
		} else {
			return nil, fmt.Errorf("Unexpected line: %s", line)
		}
	}

	/* OK, so, we have a header. Let's run with it
	 * hello (2.10-1) unstable; urgency=low */

	arguments, options := partition(header, ";")
	/* Arguments: hello (2.10-1) unstable
	 * Options:   urgency=low, other=bar */

	source, remainder := partition(arguments, "(")
	version, suite := partition(remainder, ")")

	changeLog.Source = trim(source)
	changeLog.Version = trim(version)
	changeLog.Target = trim(suite)

	changeLog.Arguments = map[string]string{}

	for _, entry := range strings.Split(options, ",") {
		key, value := partition(trim(entry), "=")
		changeLog.Arguments[trim(key)] = trim(value)
	}

	var signoff string
	/* OK, we've got the header. Let's zip down. */
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(line, " ") && trim(line) != "" {
			return nil, fmt.Errorf("Error! Didn't get ending line!")
		}

		if strings.HasPrefix(line, " -- ") {
			signoff = line
			break
		}

		changeLog.Changelog = changeLog.Changelog + line
	}

	/* Right, so we have a signoff line */
	_, signoff = partition(signoff, "--")  /* Get rid of the leading " -- " */
	whom, when := partition(signoff, "  ") /* Split on the "  " */
	changeLog.ChangedBy = trim(whom)
	changeLog.When = when

	return &changeLog, nil
}

func Parse(reader io.Reader) ([]ChangelogEntry, error) {
	stream := bufio.NewReader(reader)
	ret := []ChangelogEntry{}
	for {
		entry, err := ParseOne(stream)
		if err == io.EOF {
			break
		}
		if err != nil {
			return []ChangelogEntry{}, err
		}
		ret = append(ret, *entry)
	}
	return ret, nil
}

// hello (2.10-1) unstable; urgency=low
//
//   * New upstream release.
//   * debian/patches: Drop 01-fix-i18n-of-default-message, no longer needed.
//   * debian/patches: Drop 99-config-guess-config-sub, no longer needed.
//   * debian/rules: Drop override_dh_auto_build hack, no longer needed.
//   * Standards-Version: 3.9.6 (no changes for this).
//
//  -- Santiago Vila <sanvila@debian.org>  Sun, 22 Mar 2015 11:56:00 +0100
