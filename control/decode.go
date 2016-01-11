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

	"golang.org/x/crypto/openpgp"
)

func Unmarshal(data interface{}, foo io.Reader) error {
	return nil
}

// Decoder {{{

type Decoder struct {
	paragraphReader ParagraphReader
}

// NewDecoder {{{

func NewDecoder(reader io.Reader, keyring *openpgp.EntityList) (*Decoder, error) {
	ret := Decoder{}
	pr, err := NewParagraphReader(reader, keyring)
	if err != nil {
		return nil, err
	}
	ret.paragraphReader = *pr
	return &ret, nil
}

// }}}

// Decode {{{

func (d *Decoder) Decode(into interface{}) error {
	return decode(&d.paragraphReader, reflect.ValueOf(into))
}

// Top-level decode dispatch {{{
func decode(p *ParagraphReader, into reflect.Value) error {
	if into.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Decode can only decode into a pointer!")
	}

	switch into.Elem().Type().Kind() {
	case reflect.Struct:
		paragraph, err := p.Next()
		if err != nil {
			return err
		}
		return decodeStruct(*paragraph, into)
	case reflect.Slice:
		return nil
		return decodeSlice(p, into)
	default:
		return fmt.Errorf("Can't Decode into a %s", into.Elem().Type().Name)
	}

	return nil
}

// }}}

// Top-level struct dispatch {{{

func decodeStruct(p Paragraph, into reflect.Value) error {
	/* If we have a pointer, let's follow it */
	if into.Type().Kind() == reflect.Ptr {
		return decodeStruct(p, into.Elem())
	}

	/* Store the Paragraph type for later use when checking Anonymous
	 * values. */
	paragraphType := reflect.TypeOf(Paragraph{})

	/* Right, now, we're going to decode a Paragraph into the struct */

	for i := 0; i < into.NumField(); i++ {
		field := into.Field(i)
		fieldType := into.Type().Field(i)

		/* First, let's get the name of the field as we'd index into the
		 * map[string]string. */
		paragraphKey := fieldType.Name
		if it := fieldType.Tag.Get("control"); it != "" {
			paragraphKey = it
		}

		if paragraphKey == "-" {
			/* If the key is "-", lets go ahead and skip it */
			continue
		}

		/* Now, if we have an Anonymous field, we're either going to
		 * set it to the Paragraph if it's the Paragraph Anonymous member,
		 * or, more likely, continue through */
		if fieldType.Anonymous {
			if fieldType.Type == paragraphType {
				/* Neat! Let's give the struct this data */
				field.Set(reflect.ValueOf(p))
			} else {
				/* Otherwise, we're going to avoid doing more maths on it */
				continue
			}
		}

		if value, ok := p.Values[paragraphKey]; ok {
			if err := decodeStructValue(field, fieldType, value); err != nil {
				return err
			}
			continue
		} else {
			/* XXX: Handle required fields here */
			continue
		}
	}

	return nil
}

// }}}

// set a struct field value {{{

func decodeStructValue(field reflect.Value, fieldType reflect.StructField, value string) error {
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Int:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(value))
		return nil
	case reflect.Slice:
		return decodeStructValueSlice(field, fieldType, value)
	case reflect.Struct:
		return nil
	}

	return fmt.Errorf("Unknown type of field: %s", field.Type())

}

// }}}

// set a struct field value {{{

func decodeStructValueSlice(field reflect.Value, fieldType reflect.StructField, value string) error {
	underlyingType := field.Type().Elem()

	var delim = " "
	if it := fieldType.Tag.Get("delim"); it != "" {
		delim = it
	}

	var strip = ""
	if it := fieldType.Tag.Get("strip"); it != "" {
		strip = it
	}

	value = strings.Trim(value, strip)

	for _, el := range strings.Split(value, delim) {
		el = strings.Trim(el, strip)

		targetValue := reflect.New(underlyingType)
		err := decodeStructValue(targetValue.Elem(), fieldType, el)
		if err != nil {
			return err
		}
		field.Set(reflect.Append(field, targetValue.Elem()))
	}

	return nil
}

// }}}

// Top-level slice dispatch {{{

func decodeSlice(p *ParagraphReader, into reflect.Value) error {
	return nil
}

// }}}

// }}}

// Signer {{{

func (d *Decoder) Signer() *openpgp.Entity {
	return d.paragraphReader.Signer()
}

// }}}

// }}}

// vim: foldmethod=marker
