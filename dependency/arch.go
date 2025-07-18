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

package dependency // import "pault.ag/go/debian/dependency"

import (
	"fmt"
	"reflect"
	"strings"
)

/*
 */
type Arch struct {
	ABI  string
	Libc string
	OS   string
	CPU  string
}

func ParseArchitectures(arch string) ([]Arch, error) {
	ret := []Arch{}
	arches := strings.Split(arch, " ")
	for _, el := range arches {
		el := strings.Trim(el, " \t\n\r")

		if el == "" {
			continue
		}

		arch, err := ParseArch(el)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *arch)
	}
	return ret, nil
}

func (arch *Arch) UnmarshalControl(data string) error {
	return parseArchInto(arch, data)
}

func ParseArch(arch string) (*Arch, error) {
	ret := &Arch{
		ABI:  "any",
		Libc: "any",
		OS:   "any",
		CPU:  "any",
	}
	return ret, parseArchInto(ret, arch)
}

func expandAnyAnyAny(parts []string) []string {
	if len(parts) == 0 {
		return parts
	}

	if parts[0] != "any" {
		return parts
	}

	for idx := len(parts); idx < 4; idx++ {
		parts = append([]string{"any"}, parts...)
	}

	return parts
}

func matchShortAgainstTupletable(parts []string) []string {
	var matches = strings.Split(TUPLETABLE, "\n")
nextMatch:
	for _, line := range matches {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			panic("internal error: malformed tupletable")
		}

		shortName := strings.Split(fields[1], "-")

		if len(shortName) != len(parts) {
			continue
		}

		var cpu = ""
		for idx := range parts {
			matchName := shortName[idx]
			partName := parts[idx]

			if matchName == "<cpu>" {
				cpu = partName
				continue
			}

			if matchName != partName {
				continue nextMatch
			}
		}
		return strings.Split(strings.Replace(fields[0], "<cpu>", cpu, 1), "-")
	}

	return parts
}

func archToStringWithTupletable(arch Arch) string {

	if arch == All {
		return "all"
	}

	if arch == Any {
		return "any"
	}

	rawName := []string{
		arch.ABI,
		arch.Libc,
		arch.OS,
		arch.CPU,
	}

	var matches = strings.Split(TUPLETABLE, "\n")
	for _, line := range matches {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			panic("internal error: malformed tupletable")
		}

		// try this out with our cpu and see if it matches
		longNameStr := strings.Replace(fields[0], "<cpu>", arch.CPU, 1)
		longName := strings.Split(longNameStr, "-")
		if len(longName) != 4 {
			panic("internal error: malformed tupletable long name")
		}

		if reflect.DeepEqual(longName, rawName) {
			return strings.Replace(fields[1], "<cpu>", arch.CPU, 1)
		}
	}

	return strings.Join(rawName, "-")
}

/*
 */
func parseArchInto(ret *Arch, arch string) error {

	if arch == "all" {
		*ret = Arch{
			ABI:  "all",
			Libc: "all",
			OS:   "all",
			CPU:  "all",
		}
		return nil
	}

	parts := strings.SplitN(arch, "-", 4)
	parts = expandAnyAnyAny(parts)
	parts = matchShortAgainstTupletable(parts)

	switch len(parts) {
	case 4:
		*ret = Arch{
			ABI:  parts[0],
			Libc: parts[1],
			OS:   parts[2],
			CPU:  parts[3],
		}
	case 3:
		*ret = Arch{
			ABI:  "base",
			Libc: parts[0],
			OS:   parts[1],
			CPU:  parts[2],
		}
	case 2:
		*ret = Arch{
			ABI:  "base",
			Libc: "gnu",
			OS:   parts[0],
			CPU:  parts[1],
		}
	case 1:
		*ret = Arch{
			ABI:  "base",
			Libc: "gnu",
			OS:   "linux",
			CPU:  parts[0],
		}
	case 0:
		return fmt.Errorf("parseArchInto: empty string!")
	}

	return nil
}

/*
 */
func (set *ArchSet) Matches(other *Arch) bool {
	/* If [!amd64 sparc] matches gnu-linux-any */

	if len(set.Architectures) == 0 {
		/* We're not a thing. Always true. */
		return true
	}

	not := set.Not
	for _, el := range set.Architectures {
		if el.Is(other) {
			/* For each arch; check if it matches. If it does, then
			 * return true (unless we're negated) */
			return !not
		}
	}
	/* Otherwise, let's return false (unless we're negated) */
	return not
}

/*
 */
func (arch *Arch) IsWildcard() bool {
	if arch.CPU == "all" {
		return false
	}

	if arch.ABI == "any" || arch.OS == "any" || arch.CPU == "any" || arch.Libc == "any" {
		return true
	}
	return false
}

/*
 */
func (arch *Arch) Is(other *Arch) bool {

	if arch.IsWildcard() && other.IsWildcard() {
		/* We can't compare wildcards to other wildcards. That's just
		 * insanity. We always need a concrete arch. Not even going to try. */
		return false
	} else if arch.IsWildcard() {
		/* OK, so we're a wildcard. Let's defer to the other
		 * struct to deal with this */
		return other.Is(arch)
	}

	if (arch.CPU == other.CPU || (arch.CPU != "all" && other.CPU == "any")) &&
		(arch.OS == other.OS || other.OS == "any") &&
		(arch.Libc == other.Libc || other.Libc == "any") &&
		(arch.ABI == other.ABI || other.ABI == "any") {

		return true
	}

	return false
}

// vim: foldmethod=marker
