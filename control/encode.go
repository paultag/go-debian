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
	"strconv"
	"strings"
)

// Marshalable {{{

type Marshalable interface {
	MarshalControl() (string, error)
}

// }}}

// ConvertToParagraph {{{

func ConvertToParagraph(incoming interface{}) (*Paragraph, error) {
	data := reflect.ValueOf(incoming)
	if data.Type().Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Can only Decode a pointer to a Struct")
	}
	return convertToParagraph(data.Elem())
}

// Top-level conversaion dispatch {{{

func convertToParagraph(data reflect.Value) (*Paragraph, error) {
	order := []string{}
	values := map[string]string{}

	if data.Type().Kind() != reflect.Struct {
		return nil, fmt.Errorf("Can only Decode a Struct")
	}

	paragraphType := reflect.TypeOf(Paragraph{})
	var foundParagraph Paragraph = Paragraph{}

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		fieldType := data.Type().Field(i)

		if fieldType.Anonymous {
			if fieldType.Type == paragraphType {
				foundParagraph = field.Interface().(Paragraph)
			}
			continue
		}

		paragraphKey := fieldType.Name
		if it := fieldType.Tag.Get("control"); it != "" {
			paragraphKey = it
		}

		if paragraphKey == "-" {
			/* If the key is "-", lets go ahead and skip it */
			continue
		}

		data, err := marshalStructValue(field, fieldType)
		if err != nil {
			return nil, err
		}

		order = append(order, paragraphKey)
		values[paragraphKey] = data
	}
	para := foundParagraph.Update(Paragraph{Order: order, Values: values})
	return &para, nil
}

// }}}

// convert a struct value {{{

func marshalStructValue(field reflect.Value, fieldType reflect.StructField) (string, error) {
	switch field.Type().Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Uint:
		return strconv.Itoa(int(field.Uint())), nil
	case reflect.Ptr:
		return marshalStructValue(field.Elem(), fieldType)
	case reflect.Slice:
		return marshalStructValueSlice(field, fieldType)
	case reflect.Struct:
		return marshalStructValueStruct(field, fieldType)
	}
	return "", fmt.Errorf("Unknown type")
}

// }}}

// convert a struct value of type struct {{{

func marshalStructValueStruct(field reflect.Value, fieldType reflect.StructField) (string, error) {
	/* Right, so, we've got a type we don't know what to do with. We should
	 * grab the method, or throw a shitfit. */
	if marshal, ok := field.Interface().(Marshalable); ok {
		return marshal.MarshalControl()
	}

	return "", fmt.Errorf(
		"Type '%s' does not implement control.Marshalable",
		field.Type().Name(),
	)
}

// }}}

// convert a struct value of type slice {{{

func marshalStructValueSlice(field reflect.Value, fieldType reflect.StructField) (string, error) {
	var delim = " "
	if it := fieldType.Tag.Get("delim"); it != "" {
		delim = it
	}
	data := []string{}

	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		if stringification, err := marshalStructValue(elem, fieldType); err != nil {
			return "", err
		} else {
			data = append(data, stringification)
		}
	}

	return strings.Join(data, delim), nil
}

// }}}

// }}}

// Marshal {{{

func Marshal(writer io.Writer, data interface{}) error {
	encoder, err := NewEncoder(writer)
	if err != nil {
		return err
	}
	return encoder.Encode(data)
}

// }}}

// Encoder {{{

type Encoder struct {
	writer io.Writer
}

// NewEncoder {{{

func NewEncoder(writer io.Writer) (*Encoder, error) {
	return &Encoder{writer: writer}, nil
}

// }}}

// Encode {{{

func (e *Encoder) Encode(incoming interface{}) error {
	data := reflect.ValueOf(incoming)
	return e.encode(data)
}

// Top-level Encode reflect dispatch {{{

func (e *Encoder) encode(data reflect.Value) error {
	if data.Type().Kind() == reflect.Ptr {
		return e.encode(data.Elem())
	}

	switch data.Type().Kind() {
	case reflect.Slice:
		return e.encodeSlice(data)
	case reflect.Struct:
		return e.encodeStruct(data)
	}
	return fmt.Errorf("Unknown type")
}

// }}}

// Encode a Slice {{{

func (e *Encoder) encodeSlice(data reflect.Value) error {
	for i := 0; i < data.Len(); i++ {
		if err := e.encodeStruct(data.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

// }}}

// Encode a Struct {{{

func (e *Encoder) encodeStruct(data reflect.Value) error {
	paragraph, err := convertToParagraph(data)
	if err != nil {
		return err
	}
	return paragraph.WriteTo(e.writer)
}

// }}}

// }}}

// }}}

// vim: foldmethod=marker
