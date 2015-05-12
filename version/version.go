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

package version

import (
	"errors"
	"strconv"
	"strings"
)

//
type Version struct {
	Native   bool
	Epoch    int
	Version  string
	Revision string
}

//
func rSplit(in string, piviot string) (*string, *string) {
	if strings.Contains(in, piviot) {
		for i := len(in) - 1; i > 0; i-- {
			if string(in[i]) == piviot {
				l := in[:i]
				r := in[i+1:]
				return &l, &r
			}
		}
	}
	return &in, nil
}

//
func Parse(in string) (*Version, error) {
	version := &Version{
		Epoch:    0,
		Version:  "",
		Revision: "",
		Native:   false,
	}

	in = strings.Trim(in, " \n\t\r")
	components := strings.SplitN(in, ":", 2)

	switch len(components) {
	case 1:
		version.Epoch = 0
		in = components[0]
	case 2:
		epoch, err := strconv.Atoi(components[0])
		if err != nil {
			return nil, err
		}
		in = components[1]
		version.Epoch = epoch
	default:
		return nil, errors.New("ohshit")
	}

	ver, debversion := rSplit(in, "-")

	version.Native = debversion == nil

	version.Version = *ver

	if !version.Native {
		version.Revision = *debversion
	}

	return version, nil
}
