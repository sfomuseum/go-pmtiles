package fs

import (
	"io"
	io_fs "io/fs"
)

type pmTilesFile struct {
	io_fs.File
	reader io.ReadCloser
	info   *pmTilesFileInfo
}

func (f *pmTilesFile) Stat() (io_fs.FileInfo, error) {
	return f.info, nil
}

func (f *pmTilesFile) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *pmTilesFile) Close() error {
	return f.reader.Close()
}
