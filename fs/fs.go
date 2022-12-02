package fs

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"gocloud.dev/blob"
	io_fs "io/fs"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

type PMTilesFS struct {
	io_fs.FS
	loop     *pmtiles.Loop
	database string
	mod_time time.Time
}

func New(ctx context.Context, tile_path string, database string) (io_fs.FS, error) {

	b, err := blob.OpenBucket(ctx, tile_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open tile path, %w", err)
	}

	defer b.Close()

	database_name := fmt.Sprintf("%s.pmtiles", database)
	r, err := b.NewReader(ctx, database_name, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", database_name, err)
	}

	defer r.Close()

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
		mod_time: r.ModTime(),
	}

	return pmtiles_fs, nil
}

func (pmtiles_fs *PMTilesFS) Open(path string) (io_fs.File, error) {

	ctx := context.Background()

	fq_path, err := url.JoinPath("/", pmtiles_fs.database, path)

	if err != nil {
		return nil, fmt.Errorf("Failed to join path, %w", err)
	}

	status_code, headers, body := pmtiles_fs.loop.Get(ctx, fq_path)

	if status_code != 200 {
		return nil, io_fs.ErrNotExist
	}

	str_len, ok := headers["Content-Length"]

	if !ok {
		return nil, fmt.Errorf("Missing content length")
	}

	len, err := strconv.ParseInt(str_len, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse content length '%s', %w", str_len, err)
	}

	bytes_r := bytes.NewReader(body)
	gzip_r, err := gzip.NewReader(bytes_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gzip reader, %w", err)
	}

	fname := filepath.Base(path)

	info := &pmTilesFileInfo{
		name:     fname,
		size:     len,
		mod_time: pmtiles_fs.mod_time,
	}

	f := &pmTilesFile{
		reader: gzip_r,
		info:   info,
	}

	return f, nil
}
