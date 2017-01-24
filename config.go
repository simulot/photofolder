package main

import (
	"os"
	"text/template"

	"github.com/pkg/errors"

	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

const defaultTemplate = `/{{.YYYY}}/{{.YYYY}}.{{.MM}}/{{.YYYY}}.{{.MM}}.{{.DD}}`

func readConfig() *appConfig {
	conf := struct {
		repository     checkedPath
		path           checkedPath
		folderTpl      myTemplate
		deletePatterns []string
		dryRun         bool
		deleteSmall    bool
	}{}

	readTemplate(&conf.folderTpl, kingpin.Flag("model", "model for path").Default(defaultTemplate).Short('m'))
	kingpin.Flag("dryrun", "show actions to be done, but doesn't touch files").Short('d').Default("true").BoolVar(&conf.dryRun)
	kingpin.Flag("delete", "to be deleted file patterns, like thumb*.* or picasa.ini").Default("Thumbs.db", "@__thumb", ".@__thumb").StringsVar(&conf.deletePatterns)
	kingpin.Flag("delete-small", "delete small image smaller than 256x256 pixels").Default("false").BoolVar(&conf.deleteSmall)
	readPath(&conf.repository, kingpin.Arg("repository", "media repository"))
	readPath(&conf.path, kingpin.Arg("path", "path to be cleaned, if empty, the whole repository is cleanned"))
	kingpin.Parse()
	if len(conf.path) == 0 {
		conf.path = conf.repository
	}
	for _, d := range conf.deletePatterns {
		_, err := filepath.Match(d, "test.tst")
		dieOnError(errors.Wrapf(err, "Delete patterns '%s'", d))
	}

	return &appConfig{
		path:              string(conf.path),
		repository:        string(conf.repository),
		folderTpl:         conf.folderTpl.Template,
		dryRun:            conf.dryRun,
		folderToBeChecked: newFolderList(),
		deletePatterns:    conf.deletePatterns,
		deleteSmall:       conf.deleteSmall,
	}
}

type myTemplate struct {
	*template.Template
}

func (myTemplate) String() string { return "" }

func (t *myTemplate) Set(s string) error {
	tpl, err := template.New("").Parse(s)
	if err != nil {
		return err
	}
	t.Template = tpl
	return nil
}

func readTemplate(tpl *myTemplate, s kingpin.Settings) {
	s.SetValue(tpl)
}

type checkedPath string

func (*checkedPath) String() string { return "" }

func (p *checkedPath) Set(s string) error {
	_, err := os.Stat(s)
	if err != nil {
		return err
	}
	*p = checkedPath(s)
	return nil
}

func readPath(p *checkedPath, s kingpin.Settings) {
	s.SetValue(p)
}
