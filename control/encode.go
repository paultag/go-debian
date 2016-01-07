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

func marshalStruct(incoming reflect.Value, output io.Writer) error {
	if incoming.Type().Kind() == reflect.Ptr {
		/* If we have a pointer, let's follow it */
		return marshalStruct(incoming.Elem(), output)
	}

	/* Check to make sure the struct is sane, and try the following things:
	 *
	 * Iterate over all fields in the struct, skipping any Anonymous
	 * members.
	 */
	for i := 0; i < incoming.NumField(); i++ {
		field := incoming.Field(i)
		fieldType := incoming.Type().Field(i)

		paragraphKey := fieldType.Name
		if it := fieldType.Tag.Get("control"); it != "" {
			paragraphKey = it
		}

		if paragraphKey == "-" || fieldType.Anonymous {
			continue
		}

		val, err := marshalStructValue(field, fieldType)
		if err != nil {
			return err
		}

		if val == "" {
			continue
		}

		output.Write([]byte(paragraphKey))
		output.Write([]byte(": "))
		output.Write([]byte(dotEncode(val)))
		output.Write([]byte("\n"))
	}
	return nil
}

func dotEncode(input string) string {
	el := strings.Replace(input, "\n", "\n ", -1)
	return strings.Replace(el, "\n \n", "\n .\n", -1)
}

func marshalStructValue(field reflect.Value, fieldType reflect.StructField) (string, error) {
	/*
	 *   - If the type's a builtin, do the right thing with it.
	 *   - If the type's a struct, decode the struct value.
	 *   - If the type's a list, decode and concat the values.
	 *
	 *     If the value is empty-string, skip that Key entirely unless
	 *     it's marked as Required.
	 */
	switch field.Type().Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int:
		return strconv.Itoa(int(field.Int())), nil
	case reflect.Slice:
		return marshalStructValueSlice(field, fieldType)
	case reflect.Struct:
		return marshalStructValueStruct(field, fieldType)
	}
	return "", fmt.Errorf("Unknown type of field: %s", field.Type())
}

func marshalStructValueSlice(field reflect.Value, fieldType reflect.StructField) (string, error) {
	vals := []string{}
	for i := 0; i < field.Len(); i++ {
		value := field.Index(i)

		str, err := marshalStructValue(value, fieldType)
		if err != nil {
			return "", err
		}
		vals = append(vals, str)
	}

	delim := " "
	if it := fieldType.Tag.Get("delim"); it != "" {
		delim = it
	}

	return strings.Join(vals, delim), nil
}

func marshalStructValueStruct(field reflect.Value, fieldType reflect.StructField) (string, error) {
	elem := field.Addr()

	if unmarshal, ok := elem.Interface().(Marshalable); ok {
		return unmarshal.MarshalControl()
	}

	return "", fmt.Errorf(
		"Type '%s' does not implement control.Marshalable",
		field.Type().Name(),
	)
}

func marshalSlice(incoming reflect.Value, output io.Writer) error {
	/* For top-level slices, where we have paragraphs delim'd by
	 * newlines */
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
