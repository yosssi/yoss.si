package main

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/staticbin"
	"github.com/yosssi/rendergold"
)

var (
	version = strings.TrimLeft(runtime.Version(), "go")
)

func main() {
	m := staticbin.Classic(Asset)
	m.Use(rendergold.Renderer(rendergold.Options{Asset: Asset}))
	m.Get("/", func(r render.Render) {
		r.HTML(
			http.StatusOK,
			"top/index",
			map[string]interface{}{
				"Version":    version,
				"Production": martini.Env == martini.Prod,
			},
		)
	})
	m.Run()
}
