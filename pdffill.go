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

// PDF form filling wrapper for PDFtk. Very crude and simple, but it does what I need.
package pdffill

import "fmt"
import "bytes"
import "io/ioutil"
import "os"
import "os/exec"
import "path/filepath"

// MergeFormData generates a temp FDF file with the provided form data, calls PDFtk, and then deletes the temp file.
func MergeFormData(form map[string]Value, pdfPath, destpath string, flatten bool) error {
	tmpDir, err := ioutil.TempDir("", "pdffill_")
	if err != nil {
		return err
	}
	defer func() {
		// Ignore error, not much we can do about it...
		os.RemoveAll(tmpDir)
	}()

	// Create and call a closure so we have an "inner defer scope"
	fdfpath := filepath.Clean(tmpDir + "/output.fdf")
	err = func() error {
		f, err := os.Create(fdfpath)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = WriteFDF(f, form)
		return err
	}()
	if err != nil {
		return err
	}

	return MergeFDF(fdfpath, pdfPath, destpath, flatten)
}

// MergeFDF will call PDFtk to merge an FDF file into a PDF file, outputting the result to another file.
func MergeFDF(fdfPath, pdfPath, destpath string, flatten bool) error {
	var err error
	fdfPath, err = filepath.Abs(fdfPath)
	if err != nil {
		return err
	}
	pdfPath, err = filepath.Abs(pdfPath)
	if err != nil {
		return err
	}
	destpath, err = filepath.Abs(destpath)
	if err != nil {
		return err
	}

	// Not the best way to do this, but oh well.
	_, err = os.Stat(fdfPath)
	if err != nil {
		return err
	}
	_, err = os.Stat(pdfPath)
	if err != nil {
		return err
	}

	// Make sure we have PDFtk
	_, err = exec.LookPath("pdftk")
	if err != nil {
		return err
	}

	command := []string{
		pdfPath,
		"fill_form", fdfPath,
		"output", destpath,
	}
	if flatten {
		command = append(command, "flatten")
	}

	// Do some crap to get any error, then run PDFtk
	var stderr bytes.Buffer
	cmd := exec.Command("pdftk", command...)
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s", bytes.TrimSpace(stderr.Bytes()))
	}
	return nil
}
