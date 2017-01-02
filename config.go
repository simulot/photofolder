package main

import (
	"os"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"
)

const defaultTemplate = `/{{.YYYY}}/{{.YYYY}}.{{.MM}}/{{.YYYY}}.{{.MM}}.{{.DD}}`

func readConfig() *appConfig {
	conf := struct {
		repository checkedPath
		path       checkedPath
		folderTpl  myTemplate
		dryRun     bool
	}{}
	readTemplate(&conf.folderTpl, kingpin.Flag("model", "model for path").Default(defaultTemplate).Short('m'))
	kingpin.Flag("dryrun", "don't touch files").Short('d').Default("true").BoolVar(&conf.dryRun)
	readPath(&conf.repository, kingpin.Arg("repository", "media repository").Required())
	readPath(&conf.path, kingpin.Arg("path", "path to be cleand "))
	kingpin.Parse()
	if len(conf.path) == 0 {
		conf.path = conf.repository
	}
	return &appConfig{
		path:              string(conf.path),
		repository:        string(conf.repository),
		folderTpl:         conf.folderTpl.Template,
		dryRun:            conf.dryRun,
		folderToBeChecked: newFolderList(),
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
