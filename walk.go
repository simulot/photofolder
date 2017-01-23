package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

// walk parses repository and emit an item on out channel for each jpg or mov
func (a *appConfig) walk() chan *entryStruct {
	out := make(chan *entryStruct, 0)
	deleteFlag := len(a.deletePatterns) > 0
	go func() {
		err := filepath.Walk(a.path, func(path string, info os.FileInfo, err error) error {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if deleteFlag {
				fname := filepath.Base(info.Name())
				for _, pattern := range a.deletePatterns {
					b, err := filepath.Match(pattern, fname)
					dieOnError(errors.Wrapf(err, "Walker can't check '%s' on file '%s'", pattern, fname))
					if b {
						log.Println("Wipe ", path)
						if !a.dryRun {
							err := os.RemoveAll(path)
							dieOnError(errors.Wrapf(err, "Wipe can't RemoveAll '%s'", path))
							return err
						}
					}
				}
			}
			if !info.IsDir() {
				switch ext {
				case ".jpg", ".mov", ".png", ".gif", ".mp4":
					atomic.AddInt64(&a.checkedFiles, 1)
					out <- &entryStruct{
						path:    a.relativeName(path),
						ext:     ext,
						info:    info,
						invalid: info.Size() == 0,
					}
				}
			} else {
				switch info.Name() {
				case ".sync", ".stversions":
					return filepath.SkipDir
				}
				// Remember to check this folder at end, for removing it if empty
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
	path          string      // Actual file path relative to repository
	rightPath     string      // Path that follows template
	ext           string      // File extension
	info          os.FileInfo // File information
	exif          *exif.Exif  // Exif
	dateTaken     time.Time   // Date of shooting
	width, height int         // Picture size
	invalid       bool        // When file is invalid
}

func newEntry(path string, info os.FileInfo) *entryStruct {
	return &entryStruct{
		path: path,
		ext:  filepath.Ext(info.Name()),
		info: info,
		exif: nil,
	}
}
