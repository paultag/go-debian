/* Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
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
 * THE SOFTWARE. */

package control

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// A Paragraph is a block of RFC2822-like key value pairs. This struct contains
// two methods to fetch values, a Map called Values, and a Slice called
// Order, which maintains the ordering as defined in the RFC2822-like block
type Paragraph struct {
	Values map[string]string
	Order  []string
}

// Given a bufio.Reader, go through and return a Paragraph.
func ParseParagraph(reader *bufio.Reader) (ret *Paragraph, ohshit error) {

	ret = &Paragraph{
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
			return nil, fmt.Errorf("The shit.")
		}
	}

	return
}
