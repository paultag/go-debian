/* {{{ Copyright © 2012 Michael Stapelberg and contributors
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 *     * Neither the name of Michael Stapelberg nor the
 *       names of contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY Michael Stapelberg ''AS IS'' AND ANY
 * EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL Michael Stapelberg BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. }}} */

// version is a pure-go implementation of dpkg version string functions
// (parsing, comparison) which is compatible with dpkg(1).
package version // import "pault.ag/go/debian/version"

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"bytes"
)

// Slice is a slice versions, satisfying sort.Interface
type Slice []Version

func (a Slice) Len() int {
	return len(a)
}

func (a Slice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Slice) Less(i, j int) bool {
	return Compare(a[i], a[j]) < 0
}

type Version struct {
	Epoch    uint
	Version  string
	Revision string
}

func (v *Version) Empty() bool {
	return v.Epoch == 0 && v.Version == "" && v.Revision == ""
}

func (v *Version) IsNative() bool {
	return len(v.Revision) == 0
}

func (version *Version) UnmarshalControl(data string) error {
	return parseInto(version, data)
}

func (version Version) MarshalControl() (string, error) {
	return version.String(), nil
}

func (v Version) String() string {
	var result string
	if v.Epoch > 0 {
		result = strconv.Itoa(int(v.Epoch)) + ":" + v.Version
	} else {
		result = v.Version
	}
	if len(v.Revision) > 0 {
		result += "-" + v.Revision
	}
	return result
}

const fromChars = "+-.:" //ASCII sorted chars that can occur in version
const tildeTo = '/'
const nonAlphaAfter = 'z'
const numberDigitsZero = 'a'
const endOfString = '='
const versionEnd = '.'

// According to debian policy:
// Comparison starts with (possibly empty) nonnumber, then numbers and
// nonnumbers alter. At the end is number (possibly 0). Each number is
// stripped its leading zeroes. Empty number (at the end) is equal to 0.
//
// We replace characters in nonnumbers so that they comply with debian
// ordering. We add character to the end of nonnumber, this char is
// between 'A' and tilde replacement character.
// We strip leading zeroes from each number. We prefix the number
// by lowercase letter denoting number of digits of the number.
// (i.e. "a" is 0, "c12" is 12)
//
// 1.) If two numbers differ, then their representations differ before
//     reaching end of shorter representation.
// 2.) If two non-numbers differ, then their representation differ before
//     reaching end of shorter representation.
// These two facts gives, that if two (number,non-number) pairs differ,
// then we get the difference before end of the representation of the
// shorter one.
//
// 3.) Version or revision consits of non-number,number pairs (at least one).
//     First pair may have empty non-number part. Last pair may have ommited
//     zero.
// 4.) Start of non-number representation is greater than end of version or
//     revision sequence.
// This gives that if two upstream-versions or revisions differ, then the
// difference is found before the end of the shorter one.
//
func comparableVerRev(v string, w *bytes.Buffer) {
	i := 0
	for {
		for i < len(v) && !cisdigit(rune(v[i])) {
			if cisalpha(rune(v[i])) {
				w.WriteByte(v[i])
			} else {
				fr := strings.IndexRune(fromChars, rune(v[i])) + 1
				if (fr < 1) {
					w.WriteByte(tildeTo) //Tilde replacement
				} else {
					w.WriteByte(byte(nonAlphaAfter)+byte(fr)) // more than 'z'
				}
			}
			i++
		}
		w.WriteByte(endOfString) //End of string (more than Tilde replacement, less than 'A')
		bufPos := w.Len()
		w.WriteByte('$') //To be replaced later
		for i < len(v) && v[i] == '0' { i++ }
		numStart := i
		for i < len(v) && cisdigit(rune(v[i])) {
			w.WriteByte(v[i])
			i++
		}
		w.Bytes()[bufPos] = numberDigitsZero + byte(i - numStart)
		if i >= len(v) { break }
	}
}

// Comparable string. To be used as precomputed value for fast
// comparisons. E.g., in SQL table.
func (v Version) ComparableString() string {
	var buffer bytes.Buffer
	// Epoch is encoded the same way as number part in comparableVerRev()
	if v.Epoch != 0 {
		epoch := strconv.Itoa(int(v.Epoch))
		buffer.WriteByte(numberDigitsZero + byte(len(epoch)))
		buffer.WriteString(epoch)
	} else {
		buffer.WriteByte(numberDigitsZero)
	}
	buffer.WriteByte(versionEnd)
	comparableVerRev(v.Version, &buffer)
	buffer.WriteByte(versionEnd)
	comparableVerRev(v.Revision, &buffer)
	buffer.WriteByte(versionEnd)
	return buffer.String()
}

func comparableVerRevToStr(s string, pos *int) (ret string, err error) {
	const (
		STATE_STRING = iota
		STATE_BEFORE_STRING
		STATE_NUMBER
		STATE_END
	)
	var buffer bytes.Buffer
	state := STATE_STRING
	for len(s) > *pos && state != STATE_END {
		c := s[*pos]
		switch {
			case state == STATE_NUMBER:
				numLen := int(c) - numberDigitsZero
				if len(s) < *pos + 1 + numLen {
					err = fmt.Errorf("Unexpected end of version or revision.")
					return
				}
				buffer.WriteString(s[*pos+1:*pos+1+numLen])
				*pos += numLen
				state = STATE_BEFORE_STRING
			case state != STATE_BEFORE_STRING && state != STATE_STRING:
				err = fmt.Errorf("Unexpected state.")
				return
			case c == versionEnd:
				if state == STATE_BEFORE_STRING {
					state = STATE_END
				} else {
					err = fmt.Errorf("Unexpected end of version or revision.")
					return
				}
			case c == endOfString:
				state = STATE_NUMBER
			case c == tildeTo:
				buffer.WriteByte('~')
				state = STATE_STRING
			case c > nonAlphaAfter && int(c) <= nonAlphaAfter + len(fromChars):
				buffer.WriteByte(fromChars[c-nonAlphaAfter-1])
				state = STATE_STRING
			case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z'):
				buffer.WriteByte(c)
				state = STATE_STRING
			default:
				err = fmt.Errorf("Unexpected character.")
				return
		}
		*pos++
	}
	if state != STATE_END {
		err = fmt.Errorf("Unexpected end of version or revision.")
		return
	}
	*pos ++
	ret = buffer.String()
	return
}

func ComparableStringToVersion(s string) (ret Version, err error) {
	if len(s) < 4 {
		err = fmt.Errorf("Comparable version too short.")
		return
	}
	pos := int(s[0]) - numberDigitsZero + 1
	if len(s) < pos+2 {
		err = fmt.Errorf("Unexpected end of comparable version.")
		return
	}
	if pos > 1 {
		var epoch int64
		epoch, err = strconv.ParseInt(s[1:pos], 10, 32)
		if err != nil { return }
		ret.Epoch = uint(epoch)
	} else {
		ret.Epoch = 0
	}
	if s[pos] != versionEnd {
		err = fmt.Errorf("Expected version separator after epoch.")
		return
	}
	pos++
	ret.Version, err = comparableVerRevToStr(s, &pos)
	if err != nil { return }
	ret.Revision, err = comparableVerRevToStr(s, &pos)
	if err != nil { return }
	if pos < len(s) {
		err = fmt.Errorf("Extra characters after revision.")
		return
	}
	return
}

func cisdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func cisalpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func order(r rune) int {
	if cisdigit(r) {
		return 0
	}
	if cisalpha(r) {
		return int(r)
	}
	if r == '~' {
		return -1
	}
	if int(r) != 0 {
		return int(r) + 256
	}
	return 0
}

func verrevcmp(a string, b string) int {
	i := 0
	j := 0
	for i < len(a) || j < len(b) {
		var first_diff int
		for (i < len(a) && !cisdigit(rune(a[i]))) ||
			(j < len(b) && !cisdigit(rune(b[j]))) {
			ac := 0
			if i < len(a) {
				ac = order(rune(a[i]))
			}
			bc := 0
			if j < len(b) {
				bc = order(rune(b[j]))
			}
			if ac != bc {
				return ac - bc
			}
			i++
			j++
		}

		for i < len(a) && a[i] == '0' {
			i++
		}
		for j < len(b) && b[j] == '0' {
			j++
		}

		for i < len(a) && cisdigit(rune(a[i])) && j < len(b) && cisdigit(rune(b[j])) {
			if first_diff == 0 {
				first_diff = int(rune(a[i]) - rune(b[j]))
			}
			i++
			j++
		}

		if i < len(a) && cisdigit(rune(a[i])) {
			return 1
		}
		if j < len(b) && cisdigit(rune(b[j])) {
			return -1
		}
		if first_diff != 0 {
			return first_diff
		}
	}
	return 0
}

// Compare compares the two provided Debian versions. It returns 0 if a and b
// are equal, a value < 0 if a is smaller than b and a value > 0 if a is
// greater than b.
func Compare(a Version, b Version) int {
	if a.Epoch > b.Epoch {
		return 1
	}
	if a.Epoch < b.Epoch {
		return -1
	}

	rc := verrevcmp(a.Version, b.Version)
	if rc != 0 {
		return rc
	}

	return verrevcmp(a.Revision, b.Revision)
}

// Parse returns a Version struct filled with the epoch, version and revision
// specified in input. It verifies the version string as a whole, just like
// dpkg(1), and even returns roughly the same error messages.
func Parse(input string) (Version, error) {
	result := Version{}
	return result, parseInto(&result, input)
}

func validateVerRev(v string, name string, chars string) error {
	consecutive_digits := 0
	if strings.IndexFunc(v, func(c rune) bool {
		if cisdigit(c) {
			consecutive_digits += 1
			// ComparableString would fail...
			return consecutive_digits > 25 
		} else {
			consecutive_digits = 0
			return !cisalpha(c) && strings.IndexRune(chars, c) < 0
		}
	}) != -1 {
		if consecutive_digits > 0 {
			return fmt.Errorf("too big number in " + name)
		} else {
			return fmt.Errorf("invalid character in " + name)
		}
	}
	return nil
}

func parseInto(result *Version, input string) error {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return fmt.Errorf("version string is empty")
	}

	if strings.IndexFunc(trimmed, unicode.IsSpace) != -1 {
		return fmt.Errorf("version string has embedded spaces")
	}

	colon := strings.Index(trimmed, ":")
	if (colon >= 9) {
		return fmt.Errorf("epoch too big or not an integer")
	}
	if colon != -1 {
		epoch, err := strconv.ParseInt(trimmed[:colon], 10, 64)
		if err != nil {
			return fmt.Errorf("epoch: %v", err)
		}
		if epoch < 0 {
			return fmt.Errorf("epoch in version is negative")
		}
		result.Epoch = uint(epoch)
	}

	result.Version = trimmed[colon+1:]
	if hyphen := strings.LastIndex(result.Version, "-"); hyphen != -1 {
		result.Revision = result.Version[hyphen+1:]
		result.Version = result.Version[:hyphen]
	}

	/*
	// Version should start with number. Thus package with empty version
	// or version starting by nonnumber would not be accepted to debian.
	// But unofficial and personal repositories may have different rules.

	if len(result.Version) == 0 {
		return fmt.Errorf("empty version number")
	}

	if !unicode.IsDigit(rune(result.Version[0])) {
		return fmt.Errorf("version number does not start with digit")
	}
	*/

	err := validateVerRev(result.Version, "version number", ".-+~:")
	if err != nil { return err }

	err = validateVerRev(result.Revision, "revision number", ".+~")
	if err != nil { return err }

	return nil
}

// vim:ts=4:sw=4:noexpandtab foldmethod=marker
