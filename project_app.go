package project

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed view/project.tpl
var projectTemplate string

type App struct {
	Trace        bool
	Port         int    `name :"port" usage:"port to listen to"`
	ProjectsFile string `name:"f" usage:"projects file (.yaml)"`
	Template     string `name:"template" usage:"template file"`
	projects     []*Project
}

func (t *App) Init() error {
	t.Port = 8080
	return nil
}

func (t *App) Configured() error {
	if t.ProjectsFile != "" {
		data, err := os.ReadFile(t.ProjectsFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &t.projects)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintYaml(v interface{}) {
	data, err := yaml.Marshal(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Stdout.Write(data)
	fmt.Println()
}

func (t *App) List() {
	for _, p := range t.projects {
		fmt.Println(p.Package)
	}
}

func (t *App) FindProject(pkg string) *Project {
	for _, p := range t.projects {
		if pkg == p.Package {
			return p
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
	pkg := host + url.Path
	fmt.Println(pkg)
	project := t.FindProject(pkg)
	if project == nil {
		return errors.New("no such package: " + pkg)
	}
	var tpl *template.Template
	if t.Template != "" {
		tpl, err = template.ParseFiles(t.Template)
	} else {
		tpl = template.New("project")
		tpl, err = tpl.Parse(projectTemplate)
	}
	if err != nil {
		return err
	}
	return tpl.Execute(w, project)
}

func (t *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := t.Handle(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 404)

	}
}

func (t *App) Server() error {
	addr := fmt.Sprintf(":%d", t.Port)
	fmt.Println(addr)
	return http.ListenAndServe(addr, t)
}

func (t *App) PrintTemplate() {
	fmt.Println(len(projectTemplate))
	fmt.Println(projectTemplate)
}
