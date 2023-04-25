package goget

type Config struct {
	Domains  map[string]string
	Projects []*Project
}

type Project struct {
	Package    string `yaml:"package"`
	Repository string `yaml:"repository"`
}
