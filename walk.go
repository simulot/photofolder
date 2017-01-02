package main

import (
	"os"
	"path/filepath"
	"time"

	"strings"

	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

// walk parses repository and emit an item on out channel for each jpg or mov
func (a *appConfig) walk() chan *entryStruct {
	out := make(chan *entryStruct, 0)
	go func() {
		err := filepath.Walk(a.path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(info.Name()))
				switch ext {
				case ".jpg", ".mov", ".png":
					atomic.AddInt64(&a.checkedFiles, 1)
					out <- &entryStruct{
						path:    a.relativeName(path),
						ext:     ext,
						info:    info,
						invalid: info.Size() == 0,
					}
				case ".ini":
					atomic.AddInt64(&a.checkedFiles, 1)
					out <- &entryStruct{
						path:    a.relativeName(path),
						ext:     ext,
						info:    info,
						invalid: true,
					}
				}
			} else {
				// switch info.Name() {
				// case ".@__thumb":
				// 	out <- &entryStruct{
				// 		path:    a.relativeName(path),
				// 		ext:     "",
				// 		info:    info,
				// 		invalid: true,
				// 	}
				// 	return filepath.SkipDir
				// }

				a.folderToBeChecked.add(path)
			}
			return err
		})
		dieOnError(errors.Wrapf(err, "Can't parse folder '%s'", a.repository))
		close(out)
	}()

	return out
}

type entryStruct struct {
	path      string      // Actual file path relative to repository
	rightPath string      // Path that follows template
	ext       string      // File extension
	info      os.FileInfo // File information
	exif      *exif.Exif  // Exif
	dateTaken time.Time   // Date of shooting
	invalid   bool        // When file is invalid
}

func newEntry(path string, info os.FileInfo) *entryStruct {
	return &entryStruct{
		path: path,
		ext:  filepath.Ext(info.Name()),
		info: info,
		exif: nil,
	}
}
