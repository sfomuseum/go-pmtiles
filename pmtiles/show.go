package pmtiles

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	// "github.com/dustin/go-humanize"
	"io"
	"log"
	"os"
)

// Show prints detailed information about an archive.
func Show(_ *log.Logger, output io.Writer, bucketURL string, key string, showHeaderJsonOnly bool, showMetadataOnly bool, showTilejson bool, publicURL string, showTile bool, z int, x int, y int) error {
	ctx := context.Background()

	bucketURL, key, err := NormalizeBucketKey(bucketURL, "", key)

	if err != nil {
		return err
	}

	bucket, err := OpenBucket(ctx, bucketURL, "")

	if err != nil {
		return fmt.Errorf("Failed to open bucket for %s, %w", bucketURL, err)
	}
	defer bucket.Close()

	r, err := bucket.NewRangeReader(ctx, key, 0, 16384)

	if err != nil {
		return fmt.Errorf("Failed to create range reader for %s, %w", key, err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Failed to read %s, %w", key, err)
	}
	r.Close()

	header, err := DeserializeHeader(b[0:HeaderV3LenBytes])
	if err != nil {
		// check to see if it's a V2 file
		if string(b[0:2]) == "PM" {
			specVersion := b[2]
			return fmt.Errorf("PMTiles version %d detected; please use 'pmtiles convert' to upgrade to version 3", specVersion)
		}

		return fmt.Errorf("Failed to read %s, %w", key, err)
	}

	if !showTile {
		metadataReader, err := bucket.NewRangeReader(ctx, key, int64(header.MetadataOffset), int64(header.MetadataLength))
		if err != nil {
			return fmt.Errorf("Failed to create range reader for %s, %w", key, err)
		}

		metadataBytes, err := DeserializeMetadataBytes(metadataReader, header.InternalCompression)
		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", key, err)
		}

		if showMetadataOnly && showTilejson {
			return fmt.Errorf("cannot use more than one of --header-json, --metadata, and --tilejson together")
		}

		if showHeaderJsonOnly {
			fmt.Fprintln(output, headerToStringifiedJson(header))
		} else if showMetadataOnly {
			fmt.Fprintln(output, string(metadataBytes))
		} else if showTilejson {
			if publicURL == "" {
				// Using Fprintf instead of logger here, as this message should be written to Stderr in case
				// Stdout is being redirected.
				fmt.Fprintln(os.Stderr, "no --public-url specified; using placeholder tiles URL")
			}
			tilejsonBytes, err := CreateTileJSON(header, metadataBytes, publicURL)
			if err != nil {
				return fmt.Errorf("Failed to create tilejson for %s, %w", key, err)
			}
			fmt.Fprintln(output, string(tilejsonBytes))
		} else {
			fmt.Printf("pmtiles spec version: %d\n", header.SpecVersion)
			// fmt.Printf("total size: %s\n", humanize.Bytes(uint64(r.Size())))
			fmt.Printf("tile type: %s\n", tileTypeToString(header.TileType))
			fmt.Printf("bounds: (long: %f, lat: %f) (long: %f, lat: %f)\n", float64(header.MinLonE7)/10000000, float64(header.MinLatE7)/10000000, float64(header.MaxLonE7)/10000000, float64(header.MaxLatE7)/10000000)
			fmt.Printf("min zoom: %d\n", header.MinZoom)
			fmt.Printf("max zoom: %d\n", header.MaxZoom)
			fmt.Printf("center: (long: %f, lat: %f)\n", float64(header.CenterLonE7)/10000000, float64(header.CenterLatE7)/10000000)
			fmt.Printf("center zoom: %d\n", header.CenterZoom)
			fmt.Printf("addressed tiles count: %d\n", header.AddressedTilesCount)
			fmt.Printf("tile entries count: %d\n", header.TileEntriesCount)
			fmt.Printf("tile contents count: %d\n", header.TileContentsCount)
			fmt.Printf("clustered: %t\n", header.Clustered)
			internalCompression, _ := compressionToString(header.InternalCompression)
			fmt.Printf("internal compression: %s\n", internalCompression)
			tileCompression, _ := compressionToString(header.TileCompression)
			fmt.Printf("tile compression: %s\n", tileCompression)

			var metadataMap map[string]interface{}
			json.Unmarshal(metadataBytes, &metadataMap)
			for k, v := range metadataMap {
				switch v := v.(type) {
				case string:
					fmt.Println(k, v)
				default:
					fmt.Println(k, "<object...>")
				}
			}

			if strings.HasPrefix(bucketURL, "https://") {
				fmt.Println("web viewer: https://pmtiles.io/#url=" + url.QueryEscape(bucketURL+"/"+key))
			}
		}
	} else {
		// write the tile to stdout

		tileID := ZxyToID(uint8(z), uint32(x), uint32(y))

		dirOffset := header.RootOffset
		dirLength := header.RootLength

		for depth := 0; depth <= 3; depth++ {
			r, err := bucket.NewRangeReader(ctx, key, int64(dirOffset), int64(dirLength))
			if err != nil {
				return fmt.Errorf("Network error")
			}
			defer r.Close()
			b, err := io.ReadAll(r)
			if err != nil {
				return fmt.Errorf("I/O Error")
			}
			directory := DeserializeEntries(bytes.NewBuffer(b), header.InternalCompression)
			entry, ok := FindTile(directory, tileID)
			if ok {
				if entry.RunLength > 0 {
					tileReader, err := bucket.NewRangeReader(ctx, key, int64(header.TileDataOffset+entry.Offset), int64(entry.Length))
					if err != nil {
						return fmt.Errorf("Network error")
					}
					defer tileReader.Close()
					tileBytes, err := io.ReadAll(tileReader)
					if err != nil {
						return fmt.Errorf("I/O Error")
					}
					output.Write(tileBytes)
					break
				}
				dirOffset = header.LeafDirectoryOffset + entry.Offset
				dirLength = uint64(entry.Length)
			} else {
				fmt.Println("Tile not found in archive.")
				return nil
			}
		}
	}
	return nil
}
