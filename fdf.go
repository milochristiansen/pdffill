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

// PDF form filling wrapper for PDFtk.
package pdffill

import "io"
import "fmt"
import "bytes"

func FormatFDF(form map[string]Value) string {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, header)

	for k, v := range form {
		v.WriteFDFField(buf, k)
	}

	fmt.Fprint(buf, footer)
	return buf.String()
}

func WriteFDF(w io.Writer, form map[string]Value) (int, error) {
	written, err := fmt.Fprint(w, header)
	if err != nil {
		return written, err
	}

	for k, v := range form {
		n, err := v.WriteFDFField(w, k)
		written += n
		if err != nil {
			return written, err
		}
	}

	n, err := fmt.Fprint(w, footer)
	written += n
	return written, err
}

const header = `%FDF-1.2
%,,oe"
1 0 obj
<<
/FDF << /Fields [`

const footer = `]
>>
>>
endobj
trailer
<<
/Root 1 0 R
>>
%%EOF`
