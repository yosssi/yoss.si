package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/drone/routes"
	"github.com/yosssi/gold"
	"github.com/yosssi/gologger"
	"github.com/yosssi/yosssi/consts"
)

var (
	loggerConfig map[string]string
	serverConfig map[string]string
	logger       gologger.Logger
	version      = strings.TrimLeft(runtime.Version(), "go")
	templates    = make(map[string]*template.Template)
	g            *gold.Generator
)

func init() {
	loadJSON()
	setLogger()
	g = gold.NewGenerator(isProduction())
}

func main() {
	setHandle()
	serve()
}

func loadJSON() {
	setJSON(consts.ServerJsonPath, &serverConfig)
	setJSON(consts.LoggerJsonPath, &loggerConfig)
}

func setLogger() {
	logger = gologger.GetLogger(loggerConfig)
}

func setHandle() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	mux := routes.New()
	if isProduction() {
		mux.Static("/public", pwd)
	} else {
		mux.Static("/", pwd)
	}

	mux.Filter(logRequest)

	mux.Get("/", top)
	http.Handle("/", mux)
}

func serve() {
	logger.Info("Listening on port", serverConfig["Port"])
	err := http.ListenAndServe(":"+serverConfig["Port"], nil)
	if err != nil {
		panic(err)
	}
}

func setJSON(path string, config *map[string]string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(strings.NewReader(string(bytes)))
	err = dec.Decode(config)
	if err != nil {
		panic(err)
	}
}

func isProduction() bool {
	return !isDebug()
}

func isDebug() bool {
	return serverConfig["Debug"] == "true"
}

func top(w http.ResponseWriter, r *http.Request) {
	render(w, "./views/top.gold", map[string]interface{}{
		"IsProduction": isProduction(),
		"Version":      version,
	})
}

func render(w http.ResponseWriter, file string, data interface{}) {
	tpl, err := g.ParseFile(file)
	if err != nil {
		handleError(w, err)
	}
	if err := tpl.Execute(w, data); err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	logger.Error(err.Error())
	http.Error(w, consts.ErrorMessageInternalServerError, http.StatusInternalServerError)
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	logger.Info("Request:", r.URL)
}
