package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type appConfig struct {
	repository                   string // Folder to be cleaned
	path                         string
	folderTpl                    *template.Template
	dryRun                       bool
	checkedFiles, processedFiles int64
	folderToBeChecked            *folderList
}

func main() {
	app := readConfig()
	app.run()

}

func (a *appConfig) run() {
	a.process(a.readDetails(a.walk()))
	a.clean()
}

func dieOnError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (a *appConfig) fullName(relative string) string {
	return filepath.Join(a.repository, relative)
}
func (a *appConfig) relativeName(full string) string {
	return strings.TrimPrefix(full, a.repository)
}

func timedOutFunc(d time.Duration, fn func()) error {
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(d):
		return errors.New("Time out")
	}
}
