package goget

type Config struct {
	Domains map[string]string `yaml:"domains"`
	Modules []*Module         `yaml:"modules"`
}

type Module struct {
	Path       string `yaml:"package"`
	Repository string `yaml:"repository"`
}
