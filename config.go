package goget

type Config struct {
	Domains map[string]string `json:"domains"`
	Modules []*Module         `json:"modules"`
}

type Module struct {
	Path       string `json:"package"`
	Repository string `json:"repository"`
}
