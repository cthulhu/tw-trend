package store

import (
	"io"
	"os"
)

type MultiFile struct {
	currentFile int
	files       []*os.File
}

func (mf *MultiFile) Close() error {
	for _, file := range mf.files {
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (mf *MultiFile) Read(p []byte) (n int, err error) {
	n, err = mf.files[mf.currentFile].Read(p)
	if err == io.EOF && mf.currentFile < len(mf.files)-1 {
		mf.currentFile++
		err = nil
	}
	return
}
