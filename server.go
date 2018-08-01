package codeviewer

import (
	"log"
	"errors"
	"net/http"
	"html/template"
	"io/ioutil"
	"io"
	"strings"
	"os"
	"path/filepath"
	"fmt"

	"gopkg.in/urfave/cli.v1"
	"github.com/julienschmidt/httprouter"
)

// CmdServe serves a web server to look at the code
var CmdServe = cli.Command{
	Name: "serve",
	Usage: "Serve the website",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "file, f",
			Value: "",
			Usage: "File to display",
		},
		cli.StringFlag{
			Name: "lang",
			Value: "",
			Usage: "Language to highlight for (hljs class)", // TODO can also be set in overlay
		},
		cli.StringFlag{
			Name: "addr",
			Value: ":8080",
			Usage: "Adress to listen on",
		},
	},
	Action: serve,
}

var (
	codefile string
	td templateData
)


// given to the template
type templateData struct {
	Filename string
	Theme string
	Language string
	Content string
	ThemeList map[string]string
}

// serve the codepage
func serve(c *cli.Context) error {
	// validation
	if c.String("file") == "" {
		return errors.New("No input file given")
	}

	l := c.String("lang")
	if l == "" {
		l = strings.TrimPrefix(filepath.Ext(c.String("file")), ".")
	}

	codefile = c.String("file")

	// preparing templateData
	td = templateData{
		Content: "No code here",
		Filename: codefile,
		Language: l,
		Theme: getLastTheme(),
		ThemeList: getThemeList(),
	}

	// serving
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/highlight.min.js", hljs)
	r.GET("/style/:name", style)
	r.GET("/lang/:name", lang)	
	r.GET("/set-style/:name", setStyle) // used by ajax call

	log.Println("listening on " + c.String("addr"))
	http.ListenAndServe(c.String("addr"), r)
	return nil
}

// display code
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tmpl := template.New("")
	tmpl.Parse(html)

	// reread content on each request
    data, err := ioutil.ReadFile(codefile)
	if err != nil {
		log.Fatal(err)
	}

	td.Content = string(data);
	td.Theme = getLastTheme();

	err = tmpl.Execute(w, td)
	if err != nil {
		log.Fatal(err)
	}
}

// serve style
func style(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/css")
	fileToWriter(w, filepath.Join(ConfigDir, StyleDir, p.ByName("name")))
}

// serve language
func lang(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/javascript")
    fileToWriter(w, filepath.Join(ConfigDir, LangDir, p.ByName("name")))
}

// serve language
func hljs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/javascript")
    fileToWriter(w, filepath.Join(ConfigDir, "highlight.min.js"))
}

// set current theme
func setStyle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := ioutil.WriteFile(filepath.Join(ConfigDir, "last-style"), []byte(p.ByName("name")), 0777)
	if err != nil {
		log.Println("cannot set last-style")
	}
}

// open a file and write it to the writer
func fileToWriter(w io.Writer, p string) {
	f, err := os.Open(p)
	if err != nil {
		fmt.Fprint(w, "Not found")
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		fmt.Fprint(w, "Not found")
		return
	}
}

// returns all theme names
func getThemeList() map[string]string {
	l := make(map[string]string)

	files, err := ioutil.ReadDir(filepath.Join(ConfigDir, StyleDir))
	if err != nil {
		return l
	}

	// makes a name prettier
	prettyName := func(n string) string {
		n = strings.Replace(n, "-", " ", -1)
		n = strings.Replace(n, "_", " ", -1)
		n = strings.Replace(n, ".", " ", -1)

		parts := strings.Split(n, " ")
		n = ""
		for i, _ := range parts {
			parts[i] = strings.Title(parts[i])
		}

		return strings.Join(parts, " ")
	}

	var name string
	for _, fi := range files {
		// only theme files
		if !strings.HasSuffix(fi.Name(), ".min.css") {
			continue
		}

		name = strings.TrimSuffix(fi.Name(), ".min.css")
		
		l[name] = prettyName(name)
	}

	return l
}

func getLastTheme() string {
	data, err := ioutil.ReadFile(filepath.Join(ConfigDir, "last-style"))
	if err != nil {
		log.Println(err)
		return "atelier-estuary-light"
	}

	return string(data)
}