package goget

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"melato.org/goget/util"
)

//go:embed view/project.tpl
var projectTemplate string

type App struct {
	Trace        bool
	Port         int    `name :"port" usage:"port to listen to"`
	ProjectsFile string `name:"f" usage:"projects file (.yaml)"`
	Template     string `name:"template" usage:"template file"`
	modTime      time.Time
	projects     map[string]*Project
	queue     *util.Get[string, *Project]
}

func (t *App) Init() error {
	t.Port = 8080
	return nil
}

func (t *App) LoadProjects() error {
	st, err := os.Stat(t.ProjectsFile)
	if err != nil {
		return err
	}
	modTime := st.ModTime()
	if !modTime.After(t.modTime) {
		return nil
	}
	data, err := os.ReadFile(t.ProjectsFile)
	if err != nil {
		return err
	}

	var projects []*Project
	err = yaml.Unmarshal(data, &projects)
	if err != nil {
		return err
	}
	t.projects = make(map[string]*Project)
	for _, p := range projects {
		t.projects[p.Package] = p
	}
	return nil
}

func (t *App) Configured() error {
	t.queue = util.NewGet(t.FindProject)
	err := t.LoadProjects()
	return err
}

func (t *App) List() {
	for _, p := range t.projects {
		fmt.Println(p.Package)
	}
}

func (t *App) FindProject(pkg string) *Project {
	if t.LoadProjects() != nil {
		return nil
	}
	return t.projects[pkg]
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
	pkg = strings.TrimSuffix(pkg, "/")
	fmt.Println(pkg)
	project := t.queue.Get(pkg)
	if project == nil {
		return fmt.Errorf("no such package: %s", pkg)
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
		http.Error(w, "", 404)

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
