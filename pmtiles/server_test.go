package pmtiles

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegex(t *testing.T) {
	ok, key, z, x, y, ext := parse_tile_path("/foo/0/0/0")
	assert.False(t, ok)
	ok, key, z, x, y, ext = parse_tile_path("/foo/0/0/0.mvt")
	assert.True(t, ok)
	assert.Equal(t, key, "foo")
	assert.Equal(t, z, uint8(0))
	assert.Equal(t, x, uint32(0))
	assert.Equal(t, y, uint32(0))
	assert.Equal(t, ext, "mvt")
	ok, key, z, x, y, ext = parse_tile_path("/foo/bar/0/0/0.mvt")
	assert.True(t, ok)
	assert.Equal(t, key, "foo/bar")
	assert.Equal(t, z, uint8(0))
	assert.Equal(t, x, uint32(0))
	assert.Equal(t, y, uint32(0))
	assert.Equal(t, ext, "mvt")
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html
	ok, key, z, x, y, ext = parse_tile_path("/!-_.*'()/0/0/0.mvt")
	assert.True(t, ok)
	assert.Equal(t, key, "!-_.*'()")
	assert.Equal(t, z, uint8(0))
	assert.Equal(t, x, uint32(0))
	assert.Equal(t, y, uint32(0))
	assert.Equal(t, ext, "mvt")
	ok, key = parse_metadata_path("/!-_.*'()/metadata")
	assert.True(t, ok)
	assert.Equal(t, key, "!-_.*'()")
	ok, key = parse_tilejson_path("/!-_.*'().json")
	assert.True(t, ok)
	assert.Equal(t, key, "!-_.*'()")
}