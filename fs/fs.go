package fs

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/protomaps/go-pmtiles/pmtiles"
	io_fs "io/fs"
	"log"
	"net/url"
	"path/filepath"
)

type PMTilesFS struct {
	io_fs.FS
	loop     *pmtiles.Loop
	database string
}

func New(ctx context.Context, tile_path string, database string) (io_fs.FS, error) {

	logger := log.Default()
	cache_size := 64

	loop, err := pmtiles.NewLoop(tile_path, logger, cache_size, "")

	if err != nil {
		return nil, fmt.Errorf("Failed to create pmtiles loop, %w", err)
	}

	loop.Start()

	pmtiles_fs := &PMTilesFS{
		loop:     loop,
		database: database,
	}

	return pmtiles_fs, nil
}

func (pmtiles_fs *PMTilesFS) Open(path string) (io_fs.File, error) {

	ctx := context.Background()

	fq_path, err := url.JoinPath("/", pmtiles_fs.database, path)

	if err != nil {
		return nil, fmt.Errorf("Failed to join path, %w", err)
	}

	fmt.Println(fq_path)
	fmt.Println("GET")

	status_code, _, body := pmtiles_fs.loop.Get(ctx, fq_path)

	if status_code != 200 {
		return nil, io_fs.ErrNotExist
	}

	bytes_r := bytes.NewReader(body)
	gzip_r, err := gzip.NewReader(bytes_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gzip reader, %w", err)
	}

	fname := filepath.Base(path)

	info := &PMTilesFileInfo{
		name: fname,
	}

	f := &PMTilesFile{
		reader: gzip_r,
		info:   info,
	}

	return f, nil
}
