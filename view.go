// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xormcmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"text/template"

	"github.com/lunny/log"
	"xorm.io/core"
	"xorm.io/xorm"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"
)

var CmdView = &Command{
	UsageLine: "view [-s] driverName datasourceName tmplPath [generatedPath] [tableFilterReg]",
	Short:     "view a db to codes",
	Long: `
according database's tables and columns to generate codes for Go, C++ and etc.

    driverName        Database driver name, now supported four: mysql mymysql sqlite3 postgres
    datasourceName    Database connection uri, for detail infomation please visit driver's project page
    tmplPath          Template dir for generated. the default templates dir has provide 1 template
    generatedPath     This parameter is optional, if blank, the default value is models, then will
                      generated all codes in models dir
    tableFilterReg    Table name filter regexp
`,
}

func init() {
	CmdView.Run = runView
	CmdView.Flags = map[string]bool{
		"-s": false,
		"-l": false,
	}
}

type TmplView struct {
	Tables  []*core.Table
	Imports map[string]string
	Models  string
	Prefix  string
}

func dirViewsExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}

	return true
}

func runView(cmd *Command, args []string) {
	num := checkFlags(cmd.Flags, args, printReversePrompt)
	if num == -1 {
		return
	}
	args = args[num:]

	if len(args) < 3 {
		fmt.Println("params error, please see xorm help reverse")
		return
	}

	//var isMultiFile bool = true
	//if use, ok := cmd.Flags["-s"]; ok {
	//	isMultiFile = !use
	//}

	curPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	var genDir string
	var model string
	var filterPat *regexp.Regexp
	if len(args) >= 4 {
		genDir, err = filepath.Abs(args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		//[SWH|+] 经测试，path.Base不能解析windows下的“\”，需要替换为“/”
		genDir = strings.Replace(genDir, "\\", "/", -1)
		model = path.Base(genDir)

		if len(args) >= 5 {
			filterPat, err = regexp.Compile(args[4])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		model = "models"
		genDir = path.Join(curPath, model)
	}

	dir, err := filepath.Abs(args[2])
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	if !dirViewsExists(dir) {
		log.Errorf("Template %v path is not exist", dir)
		return
	}

	var langTmpl LangTmpl
	var ok bool
	var lang string = "go"
	var prefix string = "" //[SWH|+]

	cfgPath := path.Join(dir, "config")
	info, err := os.Stat(cfgPath)
	var configs map[string]string
	if err == nil && !info.IsDir() {
		configs = loadConfig(cfgPath)
		if l, ok := configs["lang"]; ok {
			lang = l
		}
		if j, ok := configs["genJson"]; ok {
			genJson, err = strconv.Atoi(j)
		}

		//[SWH|+]
		if j, ok := configs["prefix"]; ok {
			prefix = j
		}

		if j, ok := configs["ignoreColumnsJSON"]; ok {
			ignoreColumnsJSON = strings.Split(j, ",")
		}

		if j, ok := configs["created"]; ok {
			created = strings.Split(j, ",")
		}

		if j, ok := configs["updated"]; ok {
			updated = strings.Split(j, ",")
		}

		if j, ok := configs["deleted"]; ok {
			deleted = strings.Split(j, ",")
		}

	}

	if langTmpl, ok = langTmpls[lang]; !ok {
		fmt.Println("Unsupported programing language", lang)
		return
	}

	os.MkdirAll(genDir, os.ModePerm)

	supportComment = args[0] == "mysql" || args[0] == "mymysql"

	Orm, err := xorm.NewEngine(args[0], args[1])
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	tables, err := Orm.DBMetas()
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	if filterPat != nil && len(tables) > 0 {
		size := 0
		for _, t := range tables {
			if filterPat.MatchString(t.Name) {
				tables[size] = t
				size++
			}
		}
		tables = tables[:size]
	}
	f := fmt.Sprintf("%s/struct.go.tpl", dir)
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	t := template.New(f)
	t.Funcs(langTmpl.Funcs)

	tmpl, err := t.Parse(string(bs))
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	ext := ".go"
	for _, table := range tables {
		//[SWH|+]
		if prefix != "" {
			table.Name = strings.TrimPrefix(table.Name, prefix)
		}
		// imports
		tbs := []*core.Table{table}
		imports := langTmpl.GenImports(tbs)

		w, err := os.Create(path.Join(genDir, table.Name+ext))
		if err != nil {
			log.Errorf("%v", err)
			return
		}

		newbytes := bytes.NewBufferString("")

		t := &TmplView{Tables: tbs, Imports: imports, Models: model, Prefix: prefix}
		err = tmpl.Execute(newbytes, t)
		if err != nil {
			debug.PrintStack()
			log.Errorf("%v", err)
			return
		}

		tplcontent, err := ioutil.ReadAll(newbytes)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		var source string
		if langTmpl.Formater != nil {
			source, err = langTmpl.Formater(string(tplcontent))
			if err != nil {
				log.Errorf("%v-%v", err, string(tplcontent))
				source = string(tplcontent)
			}
		} else {
			source = string(tplcontent)
		}
		w.WriteString(source)
		w.Close()
	}

}
