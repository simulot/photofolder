package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/simulot/photofolder/mov"
)

type entryProcessor struct {
	a *appConfig
}

// run emits items when desired path differs from actual path
func (a *appConfig) readDetails(in chan *entryStruct) chan *entryStruct {

	out := make(chan *entryStruct, 0)
	for i := 0; i < 8; i++ {
		go func() {
			for e := range in {
				e.readInfo(a)
				if e.toBeProcessed(a) {
					out <- e
				}
			}
		}()
	}
	return out
}

func (e *entryStruct) readInfo(a *appConfig) {
	switch {
	case e.invalid:
		e.exif = &exif.Exif{}
		e.dateTaken = e.info.ModTime()
		return
	case e.ext == ".jpg":
		e.readExif(a)
		return
	case e.ext == ".mov":
		e.readMoov(a)
		return
	}
}

func (e *entryStruct) readExif(a *appConfig) {
	f, err := os.Open(a.fullName(e.path))
	defer f.Close()
	dieOnError(errors.Wrapf(err, "Can't process entry '%s'", e.FNAME()))
	e.dateTaken = e.info.ModTime()
	e.exif, err = exif.Decode(f)
	if err == nil {
		e.dateTaken, err = e.exif.DateTime()
		if !exif.IsTagNotPresentError(err) {
			dieOnError(errors.Wrapf(err, "Can't readExif for '%s'", e.path))
		}
		if err != nil {
			e.dateTaken = e.info.ModTime()
		}
	}
	if a.deleteSmall {
		f.Seek(0, 0)
		config, _, err := image.DecodeConfig(f)
		if err != nil {
			log.Println(errors.Wrapf(err, "Read image config '%s'", e.path))
			e.invalid = true
		} else {
			e.height, e.width = config.Height, config.Width
		}
	}
}
func (e *entryStruct) readMoov(a *appConfig) {
	f, err := os.Open(a.fullName(e.path))
	defer f.Close()
	dieOnError(errors.Wrapf(err, "Can't process entry '%s'", e.FNAME()))
	c, err := mov.Created(f)
	if err == nil {
		e.exif = &exif.Exif{}
		e.dateTaken = c
		return
	}
	info, err := f.Stat()
	dieOnError(errors.Wrapf(err, "Can't readMoov for '%s'", e.path))
	e.dateTaken = info.ModTime()
}

const minSize = 256 * 256

func (e *entryStruct) toBeProcessed(a *appConfig) bool {
	if e.invalid {
		return true
	}
	if a.deleteSmall {
		s := e.width * e.height
		if s > 0 && s < minSize {
			log.Printf("image %s (%dx%d) is too small", e.path, e.width, e.height)
			e.invalid = true
			return true
		}
	}
	if len(e.rightPath) == 0 {
		e.getRightPath(a)
	}
	return filepath.Dir(e.path) != e.rightPath
}

func (e *entryStruct) getRightPath(a *appConfig) {
	b := bytes.NewBuffer([]byte{})
	err := a.folderTpl.Execute(b, e)
	dieOnError(errors.Wrap(err, "Can't execute template in getRightPath"))
	e.rightPath = b.String()
}

func (e *entryStruct) dateTime() (time.Time, error) {
	if !e.dateTaken.IsZero() {
		return e.dateTaken, nil
	}
	return e.info.ModTime(), nil
}
func (e *entryStruct) formatedDateTime(fmt string) (string, error) {
	d, err := e.dateTime()
	if err != nil {
		return "", err
	}
	return d.Format(fmt), nil
}

func (e *entryStruct) YYYYMMDD() (string, error) {
	return e.formatedDateTime("20060102")
}
func (e *entryStruct) YYMMDD() (string, error) {
	return e.formatedDateTime("060102")
}
func (e *entryStruct) YYYY() (string, error) {
	return e.formatedDateTime("2006")
}
func (e *entryStruct) YY() (string, error) {
	return e.formatedDateTime("06")
}
func (e *entryStruct) MM() (string, error) {
	return e.formatedDateTime("01")
}
func (e *entryStruct) DD() (string, error) {
	return e.formatedDateTime("02")
}
func (e *entryStruct) HHMNSS() (string, error) {
	return e.formatedDateTime("150405")
}
func (e *entryStruct) HH() (string, error) {
	return e.formatedDateTime("15")
}
func (e *entryStruct) MN() (string, error) {
	return e.formatedDateTime("04")
}
func (e *entryStruct) SS() (string, error) {
	return e.formatedDateTime("05")
}
func (e *entryStruct) FNAME() string {
	return filepath.Base(e.path)
}
func (e *entryStruct) UFNAME() string {
	return strings.ToUpper(filepath.Base(e.path))
}
func (e *entryStruct) LFNAME() string {
	return strings.ToLower(filepath.Base(e.path))
}

func (e *entryStruct) EXT() string {
	return filepath.Ext(filepath.Base(e.path))
}
func (e *entryStruct) UEXT() string {
	return strings.ToUpper(filepath.Ext(filepath.Base(e.path)))
}
func (e *entryStruct) LEXT() string {
	return strings.ToLower(filepath.Ext(filepath.Base(e.path)))
}
