package ziplegacy

import (
	"archive/zip"
	"context"
	"io"

	"github.com/debu99/cicd-runner/commands/helpers/archive"
	"github.com/debu99/cicd-runner/helpers/archives"
)

// extractor is a zip stream extractor.
type extractor struct {
	r    io.ReaderAt
	size int64
	dir  string
}

// NewExtractor returns a new Zip Extractor.
func NewExtractor(r io.ReaderAt, size int64, dir string) (archive.Extractor, error) {
	return &extractor{r: r, size: size, dir: dir}, nil
}

// Extract extracts files from the reader to the directory passed to
// NewZipExtractor.
func (e *extractor) Extract(ctx context.Context) error {
	zr, err := zip.NewReader(e.r, e.size)
	if err != nil {
		return err
	}

	return archives.ExtractZipArchive(zr)
}
