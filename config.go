package goget

type Config struct {
	Domains  map[string]string
	Modules []*Module
}

type Module struct {
	Path    string `yaml:"package"`
	Repository string `yaml:"repository"`
}
