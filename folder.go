package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type folderList struct {
	sync.RWMutex
	m map[string]int
}

func newFolderList() *folderList {
	return &folderList{
		m: make(map[string]int),
	}
}
func (f *folderList) add(p string) {
	f.Lock()
	depth := strings.Count(filepath.Dir(p), string(filepath.Separator))
	f.m[p] = depth
	f.Unlock()
}

type items []struct {
	p string
	d int
}

func (a items) Len() int           { return len(a) }
func (a items) Less(i, j int) bool { return a[i].d > a[j].d }
func (a items) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (f *folderList) listByDepth() items {
	f.RLock()
	defer f.RUnlock()

	a := items{}
	for k, v := range f.m {
		a = append(a, struct {
			p string
			d int
		}{
			k,
			v,
		})
	}
	sort.Sort(a)
	return a
}

func (a *appConfig) clean() {
	for _, p := range a.folderToBeChecked.listByDepth() {
		a.cleanPath(p.p)
	}
}

func (a *appConfig) cleanPath(p string) {
	// log.Println("check emptiness of", p)
	d, err := ioutil.ReadDir(p)
	dieOnError(errors.Wrap(err, "cleanPath"))
	if len(d) == 0 {
		log.Println("remove", p)
		if !a.dryRun {
			err := os.Remove(p)
			dieOnError(errors.Wrap(err, "cleanPath"))
		}
	}
}
