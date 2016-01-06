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
	"fmt"
	"io"
	"reflect"
)

func marshalStruct(incoming reflect.Value, output io.Writer) error {
    return nil
}

func marshalSlice(incoming reflect.Value, output io.Writer) error {
    length := incoming.Len()
    for i := 0; i < length; i++ {
        data := incoming.Index(i)
        if err := marshal(data, output); err != nil {
            return err
        }
        output.Write([]byte("\n")) /* \n\n delim */
    }
    return nil
}


func Marshal(incoming interface{}, output io.Writer) error {
    val := reflect.ValueOf(incoming)
    return marshal(val, output)
}

func marshal(incoming reflect.Value, output io.Writer) error {
	if incoming.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Ouchie! Please give me a pointer!")
	}

	switch incoming.Elem().Type().Kind() {
	case reflect.Struct:
		return marshalStruct(incoming, output)
	case reflect.Slice:
		return marshalSlice(incoming, output)
	default:
		return fmt.Errorf(
			"Ouchie! I don't know how to deal with a %s",
			incoming.Elem().Type().Name,
		)

	}
}

type Marshalable interface {
	MarshalControl() (string, error)
}
