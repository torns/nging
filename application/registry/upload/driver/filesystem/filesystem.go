/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/webx-top/echo"
)

const Name = `filesystem`

var _ upload.Uploader = &Filesystem{}

func init() {
	upload.UploaderRegister(Name, func(typ string) upload.Uploader {
		return NewFilesystem(typ)
	})
}

func NewFilesystem(typ string) *Filesystem {
	uploadPath := helper.UploadDir + typ
	return &Filesystem{
		Type:       typ,
		UploadPath: uploadPath,
	}
}

type Filesystem struct {
	Type       string
	UploadPath string
}

func (f *Filesystem) Engine() string {
	return Name
}

func (f *Filesystem) filepath(fname string) string {
	return filepath.Join(echo.Wd(), strings.TrimPrefix(f.UploadPath, `/`), fname)
}

func (f *Filesystem) Put(dstFile string, src io.Reader, size int64) (string, error) {
	file := f.filepath(dstFile)
	saveDir := filepath.Dir(file)
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	view := f.UploadPath + `/` + dstFile
	//create destination file making sure the path is writeable.
	dst, err := os.Create(file)
	if err != nil {
		return view, err
	}
	defer dst.Close()
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return view, err
	}
	return view, nil
}

func (f *Filesystem) PublicURL(dstFile string) string {
	return dstFile
}

func (f *Filesystem) Get(dstFile string) (io.ReadCloser, error) {
	return f.openFile(dstFile)
}

func (f *Filesystem) openFile(dstFile string) (*os.File, error) {
	//file := f.filepath(dstFile)
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Open(file)
}

func (f *Filesystem) Delete(dstFile string) error {
	file := filepath.Join(echo.Wd(), dstFile)
	return os.Remove(file)
}

func (f *Filesystem) DeleteDir(dstDir string) error {
	dir := filepath.Join(echo.Wd(), dstDir)
	return os.RemoveAll(dir)
}
