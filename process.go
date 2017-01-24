package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

func (a *appConfig) process(in chan *entryStruct) {
	wg := sync.WaitGroup{}
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func() {
			for job := range in {
				a.processEntry(job)
				atomic.AddInt64(&a.processedFiles, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (a *appConfig) processEntry(e *entryStruct) {
	if e.invalid {
		// Source file is empty or invalid, delete it
		a.delete(e)
		return
	}
	// Does file exists in destination?
	destName := a.fullName(filepath.Join(e.rightPath, e.FNAME()))
	destStat, err := os.Stat(destName)
	if err != nil {
		// Desitation Not found
		a.move(e, destName)
		return
	}
	if destStat.Size() == 0 {
		// Destination file is empty, override it
		a.move(e, destName)
		return
	}
	if destStat.Size() != destStat.Size() {
		// Same name but different size, preserve destination
		a.moveAndPreserveDestination(e, destName)
		return
	}
	srcMd5 := filemd5(a.fullName(e.path))
	dstMd5 := filemd5(destName)
	if srcMd5 == dstMd5 {
		// Absolute same file, delete source
		a.delete(e)
		return
	}
	a.moveAndPreserveDestination(e, destName)
}

func (a *appConfig) delete(e *entryStruct) {
	log.Println("delete", e.path)
	if a.dryRun {
		return
	}
	err := os.Remove(a.fullName(e.path))
	if err != nil {
		log.Println(errors.Wrap(err, "Can't delete entry"))
	}
}

func (a *appConfig) move(e *entryStruct, dest string) {
	log.Println("rename", e.path, "in", dest)
	if a.dryRun {
		return
	}
	src := a.fullName(e.path)
	path := filepath.Dir(dest)
	err := os.MkdirAll(path, 0777)
	dieOnError(errors.Wrap(err, "Can't move entry"))
	err = os.Rename(src, dest)
	if err != nil {
		log.Println(errors.Wrap(err, "Can't move entry"))
	}
	// a.folderToBeChecked.add(filepath.Dir(src))
}

var fnregEpx = regexp.MustCompile(`(?P<main>.*?)(?P<gcount>_\d*){0,1}(?P<ext>\..*)`)

func (a *appConfig) moveAndPreserveDestination(e *entryStruct, dest string) {
	destName := filepath.Base(dest)
	destDir := filepath.Dir(dest)
	m := fnregEpx.FindSubmatchIndex([]byte(destName))
	c := 1

	realDestName := ""
	for {
		try := []byte{}
		try = fnregEpx.ExpandString(try, fmt.Sprintf("${main}_%d${ext}", c), destName, m)
		realDestName = filepath.Join(destDir, string(try))
		_, err := os.Stat(realDestName)
		if err != nil {
			// First not found is a good name
			break
		}
		c++
	}
	a.move(e, realDestName)
}

func filemd5(path string) string {
	result := []byte{}
	f, err := os.Open(path)
	dieOnError(errors.Wrapf(err, "Can't get MD5"))
	defer f.Close()

	hash := md5.New()
	_, err = io.Copy(hash, f)
	dieOnError(errors.Wrapf(err, "Can't get MD5"))
	return fmt.Sprintf("%x", hash.Sum(result))
}
