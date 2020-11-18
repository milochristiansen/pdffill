/*
Copyright 2020 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package pdffill

import "io"
import "fmt"

// Value is a type that can format itself as a valid FDF field.
type Value interface {
	// Format the value as an FDF field and return as a string.
	FDFField(key string) string

	// Format the value as an FDF field and write it to the given writer.
	WriteFDFField(w io.Writer, key string) (int, error)
}

// String is a Value implementation for strings.
type String string

func (s String) FDFField(key string) string {
	return fmt.Sprintf("<< /T (%s) /V (%v)>>\n", key, s)
}

func (s String) WriteFDFField(w io.Writer, key string) (int, error) {
	return fmt.Fprintf(w, "<< /T (%s) /V (%v)>>\n", key, s)
}

// Int is a Value implementation for strings.
type Int int

func (i Int) FDFField(key string) string {
	return fmt.Sprintf("<< /T (%s) /V (%v)>>\n", key, i)
}

func (i Int) WriteFDFField(w io.Writer, key string) (int, error) {
	return fmt.Fprintf(w, "<< /T (%s) /V (%v)>>\n", key, i)
}

// Float is a Value implementation for floating point numbers.
type Float float64

func (f Float) FDFField(key string) string {
	return fmt.Sprintf("<< /T (%s) /V (%v)>>\n", key, f)
}

func (f Float) WriteFDFField(w io.Writer, key string) (int, error) {
	return fmt.Fprintf(w, "<< /T (%s) /V (%v)>>\n", key, f)
}

// Bool is a Value implementation for strings.
type Bool bool

func (b Bool) FDFField(key string) string {
	if b {
		return fmt.Sprintf("<< /T (%s) /V /On>>\n", key)
	}
	return fmt.Sprintf("<< /T (%s) /V />>\n", key)
}

func (b Bool) WriteFDFField(w io.Writer, key string) (int, error) {
	if b {
		return fmt.Fprintf(w, "<< /T (%s) /V /On>>\n", key)
	}
	return fmt.Fprintf(w, "<< /T (%s) /V />>\n", key)
}

// MakeValue converts an int, float64, bool, or string to a Value.
func MakeValue(v interface{}) Value {
	switch t := v.(type) {
	case string:
		return String(t)
	case int:
		return Int(t)
	case float64:
		return Float(t)
	case bool:
		return Bool(t)
	default:
		return nil
	}
}
