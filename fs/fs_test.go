package fs

import (
	"context"
	"fmt"
	_ "gocloud.dev/blob/fileblob"
	"io"
	"path/filepath"
	"testing"
)

func TestPMTilesFS(t *testing.T) {

	ctx := context.Background()

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	fs_root := fmt.Sprintf("file://%s", abs_path)
	fs_database := "sfomuseum_architecture"

	f_uri := "12/655/1585.mvt"

	fs, err := New(ctx, fs_root, fs_database)

	if err != nil {
		t.Fatalf("Failed to create FS for %s, %v", fs_database, err)
	}

	f, err := fs.Open(f_uri)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", f_uri, err)
	}

	i, err := f.Stat()

	if err != nil {
		t.Fatalf("Failed stat %s, %v", f_uri, err)
	}

	if i.Name() != "1585.mvt" {
		t.Fatalf("Unexpected name, %s", i.Name())
	}

	_, err = io.Copy(io.Discard, f)

	if err != nil {
		t.Fatalf("Failed to copy file, %v", err)
	}

	err = f.Close()

	if err != nil {
		t.Fatalf("Failed to close file, %v", err)
	}

	sz := i.Size()

	if sz != int64(74801) {
		t.Fatalf("Unexpected size, %d", sz)
	}

}
