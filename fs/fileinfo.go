package fs

import (
	io_fs "io/fs"
	"time"
)

type PMTilesFileInfo struct {
	io_fs.FileInfo
	name string
}

func (i *PMTilesFileInfo) Name() string {
	return i.name
}

func (i *PMTilesFileInfo) Size() int64 {
	return 0 // fix me
}

func (i *PMTilesFileInfo) Mode() io_fs.FileMode {
	return io_fs.ModeDevice
}

func (i *PMTilesFileInfo) ModTime() time.Time {
	return time.Now() // fix me
}

func (i *PMTilesFileInfo) IsDir() bool {
	return false
}

func (i *PMTilesFileInfo) Sys() any {
	return nil
}
