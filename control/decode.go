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
	"reflect"
	"strconv"
	"strings"
)

func decodeCustomValues(incoming reflect.Value, incomingField reflect.StructField, data string) error {
	/* Incoming is a slice */
	underlyingType := incoming.Type().Elem()

	var delim = " "
	if it := incomingField.Tag.Get("delim"); it != "" {
		delim = it
	}

	var strip = ""
	if it := incomingField.Tag.Get("strip"); it != "" {
		strip = it
	}

	if strip != "" {
		data = strings.Trim(data, strip)
	}

	for _, el := range strings.Split(data, delim) {
		if strip != "" {
			el = strings.Trim(el, strip)
		}

		targetValue := reflect.New(underlyingType)
		err := decodeValue(targetValue.Elem(), incomingField, el)
		if err != nil {
			return err
		}
		incoming.Set(reflect.Append(incoming, targetValue.Elem()))
	}
	return nil
}

func decodeCustomValue(incoming reflect.Value, incomingField reflect.StructField, data string) error {
	/* Right, so, we've got a type we don't know what to do with. We should
	 * grab the method, or throw a shitfit. */
	elem := incoming.Addr()

	if unmarshal, ok := elem.Interface().(Unmarshalable); ok {
		return unmarshal.UnmarshalControl(data)
	}

	return fmt.Errorf(
		"Type '%s' does not implement control.Unmarshalable",
		incomingField.Type.Name(),
	)
}

func decodeValue(incoming reflect.Value, incomingField reflect.StructField, data string) error {
	switch incoming.Type().Kind() {
	case reflect.String:
		incoming.SetString(data)
		return nil
	case reflect.Int:
		if data == "" {
			incoming.SetInt(0)
			return nil
		}
		value, err := strconv.Atoi(data)
		if err != nil {
			return err
		}
		incoming.SetInt(int64(value))
		return nil
	case reflect.Slice:
		return decodeCustomValues(incoming, incomingField, data)
	case reflect.Struct:
		return decodeCustomValue(incoming, incomingField, data)
	}
	return fmt.Errorf("Unknown type of field: %s", incoming.Type())
}

func decodePointer(incoming reflect.Value, data Paragraph) error {
	if incoming.Type().Kind() == reflect.Ptr {
		/* If we have a pointer, let's follow it */
		return decodePointer(incoming.Elem(), data)
	}

	for i := 0; i < incoming.NumField(); i++ {
		field := incoming.Field(i)
		fieldType := incoming.Type().Field(i)

		if field.Type().Kind() == reflect.Struct {
			err := decodePointer(field, data)
			if err != nil {
				return err
			}
		}

		paragraphKey := fieldType.Name
		if it := fieldType.Tag.Get("control"); it != "" {
			paragraphKey = it
		}

		if paragraphKey == "-" {
			continue
		}

		required := fieldType.Tag.Get("required") == "true"

		if val, ok := data.Values[paragraphKey]; ok {
			err := decodeValue(field, fieldType, val)
			if err != nil {
				return fmt.Errorf(
					"pault.ag/go/debian/control: failed to set %s: %s",
					fieldType.Name,
					err,
				)
			}
		} else if required {
			return fmt.Errorf(
				"pault.ag/go/debian/control: required field %s missing",
				fieldType.Name,
			)

		}
	}

	return nil
}

// Given a struct (or list of structs), read the io.Reader RFC822-alike
// Debian control-file stream into the struct, unpacking keys into the
// struct as needed. If a list of structs is given, unpack all RFC822
// Paragraphs into the structs.
//
// This code will attempt to unpack it into the struct based on the
// literal name of the key, compared byte-for-byte. If this is not
// OK, the struct tag `control:""` can be used to define the key to use
// in the RFC822 stream.
//
// If you're unpacking into a list of strings, you have the option of defining
// a string to split tokens on (`delim:", "`), and things to strip off each
// element (`strip:"\n\r\t "`).
//
// If you're unpacking into a struct, the struct will be walked acording to
// the rules above. If you wish to override how this writes to the nested
// struct, objects that implement the Unmarshalable interface will be
// Unmarshaled via that method call only.
//
// Structs that contain Paragraph as an Anonymous member will have that
// member populated with the parsed RFC822 block, to allow access to the
// .Values and .Order members.
func Unmarshal(incoming interface{}, data io.Reader) error {
	/* Dispatch if incoming is a slice or not */
	val := reflect.ValueOf(incoming)

	if val.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Ouchie! Please give me a pointer!")
	}

	switch val.Elem().Type().Kind() {
	case reflect.Struct:
		return unmarshalStruct(incoming, data)
	case reflect.Slice:
		return unmarshalSlice(incoming, data)
	default:
		return fmt.Errorf(
			"Ouchie! I don't know how to deal with a %s",
			val.Elem().Type().Name,
		)

	}
}

// The Unmarshalable interface defines the interface that Unmarshal will use
// to do custom unpacks into Structs.
//
// The argument passed in will be a string that contains the value of the
// RFC822 key this object relates to.
type Unmarshalable interface {
	UnmarshalControl(data string) error
}

func unmarshalSlice(incoming interface{}, data io.Reader) error {
	/* Good holy hot damn this code is ugly */
	for {
		val := reflect.ValueOf(incoming)
		flavor := val.Elem().Type().Elem()

		targetValue := reflect.New(flavor)
		target := targetValue.Interface()
		err := unmarshalStruct(target, data)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		val.Elem().Set(reflect.Append(val.Elem(), targetValue.Elem()))
	}
	return nil
}

func isParagraph(incoming reflect.Value) (int, bool) {
	paragraphType := reflect.TypeOf(Paragraph{})
	val := incoming.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Anonymous && field.Type == paragraphType {
			return i, true
		}
	}
	return -1, false
}

func unmarshalStruct(incoming interface{}, data io.Reader) error {
	reader := bufio.NewReader(data)
	para, err := ParseParagraph(reader)
	if err != nil {
		return err
	}
	if para == nil {
		return io.EOF
	}

	val := reflect.ValueOf(incoming).Elem()
	/* Before we dump it back, we should give the Paragraph back to
	 * the object */
	if index, is := isParagraph(val); is {
		/* If we're a Paragraph, let's go ahead and set the index. */
		val.Field(index).Set(reflect.ValueOf(*para))
	}

	return decodePointer(reflect.ValueOf(incoming), *para)
}

// vim: foldmethod=marker
