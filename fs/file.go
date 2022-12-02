package fs

import (
	"io"
	io_fs "io/fs"
)

type PMTilesFile struct {
	io_fs.File
	reader io.ReadCloser
	info   *PMTilesFileInfo
}

func (f *PMTilesFile) Stat() (io_fs.FileInfo, error) {
	return f.info, nil
}

func (f *PMTilesFile) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *PMTilesFile) Close() error {
	return f.reader.Close()
}
