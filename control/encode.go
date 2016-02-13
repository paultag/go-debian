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
	"reflect"
	"strconv"
	"strings"
)

func Marshal(data interface{}) (*Paragraph, error) {
	return marshalStruct(reflect.ValueOf(data))
}

type Marshalable interface {
	MarshalControl() (string, error)
}

func marshalStruct(data reflect.Value) (*Paragraph, error) {
	if data.Type().Kind() == reflect.Ptr {
		return marshalStruct(data.Elem())
	}

	order := []string{}
	values := map[string]string{}

	for i := 0; i < data.NumField(); i++ {
		field := data.Field(i)
		fieldType := data.Type().Field(i)

		if fieldType.Anonymous {
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

	return &Paragraph{
		Order:  order,
		Values: values,
	}, nil
}

func marshalStructValue(field reflect.Value, fieldType reflect.StructField) (string, error) {
	switch field.Type().Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Uint:
		return strconv.Itoa(int(field.Uint())), nil
	case reflect.Slice:
		return marshalStructValueSlice(field, fieldType)
	case reflect.Struct:
		return marshalStructValueStruct(field, fieldType)
	}
	return "", fmt.Errorf("Unknown type")
}

func marshalStructValueStruct(field reflect.Value, fieldType reflect.StructField) (string, error) {
	/* Right, so, we've got a type we don't know what to do with. We should
	 * grab the method, or throw a shitfit. */
	elem := field.Addr()

	if marshal, ok := elem.Interface().(Marshalable); ok {
		return marshal.MarshalControl()
	}

	return "", fmt.Errorf(
		"Type '%s' does not implement control.Marshalable",
		field.Type().Name(),
	)
}

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

// vim: foldmethod=marker
