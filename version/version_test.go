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

package version // import "pault.ag/go/debian/version"

import (
	"testing"
	"strings"
)

// Abbreviation for creating a new Version object.
func v(epoch uint, version string, revision string) Version {
	return Version{Epoch: epoch, Version: version, Revision: revision}
}

func sgn(x int) int {
	if x < 0 {return -1}
	if x > 0 {return 1}
	return 0
}

func comparisonTest(t *testing.T, a Version, b Version, expect int) {
	cmp := Compare(a,b)
	if sgn(cmp) != sgn(expect) {
		t.Errorf("Compare(%v, %v), got %d, want %d", a, b, cmp, expect)
	}
	reprA := a.ComparableString()
	reprB := b.ComparableString()
	cmp = strings.Compare(reprA, reprB)
	if sgn(cmp) != expect {
		t.Errorf("ComparableString-Compare(%v, %v), got %d, want %d", a, b, cmp, expect)
	}
}

func TestEpoch(t *testing.T) {
	comparisonTest(t, Version{Epoch: 1}, Version{Epoch: 2}, -1);
	comparisonTest(t, v(0, "1", "1"),    v(0, "2", "1"),    -1);
	comparisonTest(t, v(0, "1", "1"),    v(0, "1", "2"),    -1);
	comparisonTest(t, v(0, "10", "1"),   v(0, "2", "1"),     1);
}

func TestEquality(t *testing.T) {
	comparisonTest(t, v(0, "0", "0"),     v(0, "0", "0"),    0);
	comparisonTest(t, v(0, "0", "00"),    v(0, "00", "0"),   0);
	comparisonTest(t, v(0, "00", ""),     v(0, "", "00"),   0);
	comparisonTest(t, v(1, "2", "3"),     v(1, "2", "3"),    0);
}

func TestEpochDifference(t *testing.T) {
	a := v(0, "0", "0")
	b := v(1, "0", "0")
	comparisonTest(t, a, b, -1)
	comparisonTest(t, b, a,  1)
}

func TestVersionDifference(t *testing.T) {
	a := v(0, "a", "0")
	b := v(0, "b", "0")
	comparisonTest(t, a, b, -1)
	comparisonTest(t, b, a,  1)
}

func TestRevisionDifference(t *testing.T) {
	a := v(0, "0", "a")
	b := v(0, "0", "b")
	comparisonTest(t, a, b, -1)
	comparisonTest(t, b, a,  1)
}

func TestComparableStringSetup(t *testing.T) {
	for i := range(fromChars) {
		if i > 0 && fromChars[i] <= fromChars[i-1] {
			t.Errorf("fromChars must be sorted!")
		}
	}
	if len(fromChars) + int('z') > 126 {
		t.Errorf("Too much fromChars. Result may be improperly encoded in national character sets.")
	}
	if (versionEnd >= tildeTo) {
		t.Errorf("version separator must be less than tilde (smallest possible start of non-number)");
	}
	if (tildeTo >= endOfString) {
		t.Errorf("tilde must be less than end of string");
	}
	if (endOfString >= 'A') {
		t.Errorf("end of string must be less than 'A'");
	}
	if (endOfString >= numberDigitsZero) {
		t.Errorf("end of string must be less than start of number");
	}
	if ('z' > nonAlphaAfter) {
		t.Errorf("'z' must be less than non-alphanumeric characters");
	}
}

func TestCompareCodesearch(t *testing.T) {
	a := v(0, "1.8.6", "2")
	b := v(0, "1.8.6", "2.1")
	comparisonTest(t, a, b, -1)
}

func TestParseZeroVersions(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0", "")

	if a, err = Parse("0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0", err)
	}
	comparisonTest(t, a, b, 0)

	if a, err = Parse("0:0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0", err)
	}
	comparisonTest(t, a, b, 0)

	if a, err = Parse("0:0-"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(0, "0", "0")
	if a, err = Parse("0:0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-0", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(0, "0.0", "0.0")
	if a, err = Parse("0:0.0-0.0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0.0-0.0", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseEpochedVersions(t *testing.T) {
	var a Version
	var err error
	b := v(1, "0", "")

	if a, err = Parse("1:0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "1:0", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(5, "1", "")
	if a, err = Parse("5:1"); err != nil {
		t.Errorf("Parsing %q failed: %v", "5:1", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseMultipleHyphens(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0-0", "0")

	if a, err = Parse("0:0-0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-0-0", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(0, "0-0-0", "0")
	if a, err = Parse("0:0-0-0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-0-0-0", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseMultipleColons(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0:0", "0")

	if a, err = Parse("0:0:0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0:0-0", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(0, "0:0:0", "0")
	if a, err = Parse("0:0:0:0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0:0:0-0", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseMultipleHyphensAndColons(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0:0-0", "0")

	if a, err = Parse("0:0:0-0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0:0-0-0", err)
	}
	comparisonTest(t, a, b, 0)

	b = v(0, "0-0:0", "0")
	if a, err = Parse("0:0-0:0-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-0:0-0", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseValidUpstreamVersionCharacters(t *testing.T) {
	var a Version
	var err error
	b := v(0, "09azAZ.-+~:", "0")

	if a, err = Parse("0:09azAZ.-+~:-0"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:09azAZ.-+~:-0", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseValidRevisionCharacters(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0", "azAZ09.+~")

	if a, err = Parse("0:0-azAZ09.+~"); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-azAZ09.+~", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseLeadingTrailingSpaces(t *testing.T) {
	var a Version
	var err error
	b := v(0, "0", "1")

	if a, err = Parse("    0:0-1"); err != nil {
		t.Errorf("Parsing %q failed: %v", "    0:0-1", err)
	}
	comparisonTest(t, a, b, 0)

	if a, err = Parse("0:0-1     "); err != nil {
		t.Errorf("Parsing %q failed: %v", "0:0-1     ", err)
	}
	comparisonTest(t, a, b, 0)

	if a, err = Parse("      0:0-1     "); err != nil {
		t.Errorf("Parsing %q failed: %v", "      0:0-1     ", err)
	}
	comparisonTest(t, a, b, 0)
}

func TestParseEmptyVersion(t *testing.T) {
	if _, err := Parse(""); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "")
	}
	if _, err := Parse("  "); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "  ")
	}
}

/*
// First, "0:" did not pass, but "0:-0" passed. Second, upsteam version
// may be in fact empty. We require only that the version string is not empty,
// to prevent errors that occur when NULL of any form is converted to empty
// string.
func TestParseEmptyUpstreamVersionAfterEpoch(t *testing.T) {
	if _, err := Parse("0:"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "0:")
	}
}
*/

func TestParseVersionWithEmbeddedSpaces(t *testing.T) {
	if _, err := Parse("0:0 0-1"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "0:0 0-1")
	}
}

func TestParseVersionWithNegativeEpoch(t *testing.T) {
	if _, err := Parse("-1:0-1"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "-1:0-1")
	}
}

func TestParseVersionWithHugeEpoch(t *testing.T) {
	if _, err := Parse("999999999999999999999999:0-1"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "999999999999999999999999:0-1")
	}
}

func TestParseInvalidCharactersInEpoch(t *testing.T) {
	if _, err := Parse("a:0-0"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "a:0-0")
	}
	if _, err := Parse("A:0-0"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "A:0-0")
	}
}

func TestParseUpstreamVersionNotStartingWithADigit(t *testing.T) {
	if _, err := Parse("0:abc3-0"); err != nil {
		// Upstream version is recommended to start with digit.
		// But we cannot force it, when dealing with external packages. -- Ebik.
		t.Errorf("Expected an ok, but %q failed: %v", "0:abc3-0", err)
	}
}

func TestParseInvalidCharactersInUpstreamVersion(t *testing.T) {
	chars := "!#@$%&/|\\<>()[]{};,_=*^'"
	for i := 0; i < len(chars); i++ {
		verstr := "0:0" + chars[i:i+1] + "-0"
		if _, err := Parse(verstr); err == nil {
			t.Errorf("Expected an error, but %q was parsed without an error", verstr)
		}
	}
}

func TestParseInvalidCharactersInRevision(t *testing.T) {
	if _, err := Parse("0:0-0:0"); err == nil {
		t.Errorf("Expected an error, but %q was parsed without an error", "0:0-0:0")
	}
	chars := "!#@$%&/|\\<>()[]{}:;,_=*^'"
	for i := 0; i < len(chars); i++ {
		verstr := "0:0-" + chars[i:i+1]
		if _, err := Parse(verstr); err == nil {
			t.Errorf("Expected an error, but %q was parsed without an error", verstr)
		}
	}
}

// vim:ts=4:sw=4:noexpandtab foldmethod=marker
