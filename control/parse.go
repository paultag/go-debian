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

package control

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/openpgp/clearsign"
)

// func encodeValue(value string) string {
// 	ret := ""
//
// 	lines := strings.Split(value, "\n")
// 	for _, line := range lines {
// 		line = strings.Trim(line, " \t\r\n")
// 		if line == "" {
// 			line = "."
// 		}
// 		line = " " + line
// 		ret = ret + line + "\n"
// 	}
//
// 	return ret
// }

// A Paragraph is a block of RFC2822-like key value pairs. This struct contains
// two methods to fetch values, a Map called Values, and a Slice called
// Order, which maintains the ordering as defined in the RFC2822-like block
type Paragraph struct {
	Values map[string]string
	Order  []string
}

// func (para Paragraph) String() string {
// 	ret := ""
//
// 	for _, key := range para.Order {
// 		value := encodeValue(para.Values[key])
// 		ret = ret + fmt.Sprintf("%s:%s", key, value)
// 	}
//
// 	return ret
// }

func ParseOpenPGPParagraph(reader *bufio.Reader) (*Paragraph, error) {
	els := ""
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		els = els + line
	}
	block, _ := clearsign.Decode([]byte(els))
	/**
	 * XXX: With the block, we need to validate everything.
	 *
	 * We need to hit openpgp.CheckDetachedSignature with block and
	 * a keyring. For now, it'll ignore all signature checking entirely.
	 */
	return ParseParagraph(bufio.NewReader(strings.NewReader(string(block.Bytes))))
}

// Given a bufio.Reader, go through and return a Paragraph.
func ParseParagraph(reader *bufio.Reader) (*Paragraph, error) {
	line, _ := reader.Peek(15)
	if string(line) == "-----BEGIN PGP " {
		return ParseOpenPGPParagraph(reader)
	}

	ret := &Paragraph{
		Values: map[string]string{},
		Order:  []string{},
	}

	var key = ""
	var value = ""
	var noop = " \n\r\t"

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(ret.Order) == 0 {
				return nil, nil
			}
			return ret, nil
		}
		if line == "\n" {
			break
		}
		if trimmed := strings.TrimLeft(line, noop); len(trimmed) > 0 && trimmed[0] == '#' {
			// skip comments
			continue
		}

		if line[0] == ' ' {
			line = line[1:]
			ret.Values[key] += "\n" + strings.Trim(line, noop)
			continue
		}

		els := strings.SplitN(line, ":", 2)

		switch len(els) {
		case 2:
			key = strings.Trim(els[0], noop)
			value = strings.Trim(els[1], noop)

			ret.Values[key] = value
			ret.Order = append(ret.Order, key)
			continue
		default:
			return nil, fmt.Errorf("Line %q is not 'key: val'", line)
		}
	}

	return ret, nil
}

// vim: foldmethod=marker
