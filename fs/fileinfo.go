package fs

import (
	io_fs "io/fs"
	"time"
)

type pmTilesFileInfo struct {
	io_fs.FileInfo
	name     string
	size     int64
	mod_time time.Time
}

func (i *pmTilesFileInfo) Name() string {
	return i.name
}

func (i *pmTilesFileInfo) Size() int64 {
	return i.size
}

func (i *pmTilesFileInfo) Mode() io_fs.FileMode {
	return io_fs.ModeDevice
}

func (i *pmTilesFileInfo) ModTime() time.Time {
	return i.mod_time
}

func (i *pmTilesFileInfo) IsDir() bool {
	return false
}

func (i *pmTilesFileInfo) Sys() any {
	return nil
}
