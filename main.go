package main

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/staticbin"
	"github.com/yosssi/ace"
)

var (
	version = strings.TrimLeft(runtime.Version(), "go")
)

func main() {
	m := staticbin.Classic(Asset)
	m.Get("/", func(w http.ResponseWriter) {
		tpl, err := ace.ParseFiles(
			"views/base",
			"views/top/index",
			&ace.Options{Asset: Asset, Cache: martini.Env == martini.Prod})
		if err != nil {
			panic(err)
		}

		data := map[string]interface{}{
			"Version":    version,
			"Production": martini.Env == martini.Prod,
		}

		if err := tpl.Execute(w, data); err != nil {
			panic(err)
		}
	})
	m.Run()
}
