package processor

import (
	"compress/gzip"
	"io"
)

// Processor is an interface that wraps an io.Writer and can be closed
type Processor interface {
	io.WriteCloser
}

// GzipProcessor wraps an io.Writer with gzip compression
type GzipProcessor struct {
	gw *gzip.Writer
}

// NewGzipProcessor creates a new GzipProcessor
func NewGzipProcessor(w io.Writer) *GzipProcessor {
	return &GzipProcessor{
		gw: gzip.NewWriter(w),
	}
}

// Write writes compressed data to the underlying writer
func (g *GzipProcessor) Write(p []byte) (n int, err error) {
	return g.gw.Write(p)
}

// Close closes the gzip writer, flushing any unwritten data
func (g *GzipProcessor) Close() error {
	return g.gw.Close()
}
