package fs

import (
	"bytes"
	"context"
	"compress/gzip"
	"fmt"
	"github.com/protomaps/go-pmtiles/pmtiles"
	io_fs "io/fs"
	"log"
	"path/filepath"
)

type PMTilesFS struct {
	io_fs.FS
	loop *pmtiles.Loop
}

func New(tile_path string) (io_fs.FS, error) {

	logger := log.Default()
	cache_size := 64

	loop, err := pmtiles.NewLoop(tile_path, logger, cache_size, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create pmtiles loop, %w", err)
	}

	pmtiles_fs := &PMTilesFS{
		loop: loop,
	}

	return pmtiles_fs, nil
}

func (pmtiles_fs *PMTilesFS) Open(path string) (io_fs.File, error) {

	ctx := context.Background()
	
	status_code, _, body := pmtiles_fs.loop.Get(ctx, path)

	if status_code != 200 {
		return nil, io_fs.ErrNotExist
	}

	br := bytes.NewReader(body)
	gr, err := gzip.NewReader(br)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gzip reader, %w", err)
	}

	fname := filepath.Base(path)

	info := &PMTilesFileInfo{
		name: fname,
	}

	f := &PMTilesFile{
		reader: gr,
		info:   info,
	}

	return f, nil
}
