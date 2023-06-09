package goget

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	text "text/template"
	"time"

	"gopkg.in/yaml.v2"
	"melato.org/goget/util"
)

//go:embed view/module.tpl
var defaultTemplate string

type App struct {
	Trace      bool
	Port       int    `name :"port" usage:"port to listen to"`
	ConfigFile string `name:"c" usage:"config file (.yaml)"`
	Template   string `name:"template" usage:"template file"`
	modTime    time.Time
	domains    map[string]*text.Template
	projects   map[string]*Module
	queue      *util.Get[specifier, *Module]
}

func (t *App) Init() error {
	t.Port = 8080
	return nil
}

func (t *App) LoadProjects() error {
	st, err := os.Stat(t.ConfigFile)
	if err != nil {
		return err
	}
	modTime := st.ModTime()
	if !modTime.After(t.modTime) {
		return nil
	}
	data, err := os.ReadFile(t.ConfigFile)
	if err != nil {
		return err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	t.domains = make(map[string]*text.Template)
	for domain, pattern := range config.Domains {
		tpl := text.New("x").Option("missingkey=error")
		tpl, err = tpl.Parse(pattern)
		if err != nil {
			return fmt.Errorf("%s: %w", domain, err)
		}
		t.domains[domain] = tpl
	}
	t.projects = make(map[string]*Module)
	for _, p := range config.Modules {
		t.projects[p.Path] = p
	}
	return nil
}

func (t *App) Configured() error {
	t.queue = util.NewGet(t.GetProject)
	err := t.LoadProjects()
	return err
}

func (t *App) List() {
	for _, p := range t.projects {
		fmt.Println(p.Path)
	}
}

type specifier struct {
	Host string
	Path string
}

func (t *App) GetProject(sp specifier) *Module {
	if t.Trace {
		fmt.Printf("host=%s path=%s\n", sp.Host, sp.Path)
	}
	if t.LoadProjects() != nil {
		return nil
	}
	pkg := sp.Host + sp.Path
	p, ok := t.projects[pkg]
	if ok {
		return p
	}
	tpl, ok := t.domains[sp.Host]
	if ok {
		model := make(map[string]string)
		model["path"] = strings.TrimPrefix(sp.Path, "/")
		var buf bytes.Buffer
		err := tpl.Execute(&buf, model)
		if err == nil {
			var p Module
			p.Path = pkg
			p.Repository = buf.String()
			return &p
		}
		if err != nil {
			fmt.Printf("%v\n", err)
			return nil
		}
	}
	return nil
}

func (t *App) host(r *http.Request) string {
	host := r.Host
	i := strings.Index(host, ":")
	if i >= 0 {
		return host[0:i]
	}
	return host
}

func (t *App) Handle(w http.ResponseWriter, r *http.Request) error {
	url, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return err
	}
	host := t.host(r)
	path := strings.TrimSuffix(url.Path, "/")
	project := t.queue.Get(specifier{Host: host, Path: path})
	if project == nil {
		return fmt.Errorf("no such package: %s%s", host, path)
	}
	if t.Trace {
		fmt.Printf("package=%s repository=%s\n", project.Path, project.Repository)
	}
	var tpl *template.Template
	if t.Template != "" {
		tpl, err = template.ParseFiles(t.Template)
	} else {
		tpl = template.New("x")
		tpl, err = tpl.Parse(defaultTemplate)
	}
	if err != nil {
		return err
	}
	tpl.Option("missingkey=error")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, project)
	if err == nil {
		buf.WriteTo(w)
	}
	return err
}

func (t *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := t.Handle(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", 404)

	}
}

func (t *App) Server() error {
	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Println(addr)
	return http.ListenAndServe(addr, t)
}

func (t *App) PrintTemplate() {
	fmt.Println(len(defaultTemplate))
	fmt.Println(defaultTemplate)
}
