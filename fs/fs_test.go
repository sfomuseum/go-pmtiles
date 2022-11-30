package fs

import (
	"context"
	"fmt"
	_ "gocloud.dev/blob/fileblob"
	"testing"
)

func TestPMTilesFS(t *testing.T) {

	ctx := context.Background()

	fs_root := "file:///usr/local/data"
	fs_database := "sfov3"

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

	fmt.Println(i.Name())
}
