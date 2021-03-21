package project

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type App struct {
	Trace        bool
	Port         int    `name :"port" usage:"port to listen to"`
	ProjectsFile string `name:"f" usage:"projects file (.yaml)"`
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

func (t *App) List() {
	for _, p := range t.projects {
		fmt.Println(p.Dir)
	}
}

func (t *App) Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
}

func (t *App) Server() error {
	r := mux.NewRouter()
	r.HandleFunc("/{p}", t.Handler)
	addr := fmt.Sprintf(":%d", t.Port)
	return http.ListenAndServe(addr, r)
}
