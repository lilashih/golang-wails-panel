package project

type Project struct {
	Config ProjectConfig

	Path    string `json:"path"`
	Running bool   `json:"running"`

	Install      func() error `json:"-"`
	Start        func() error `json:"-"`
	Stop         func() error `json:"-"`
	CheckRunning func()       `json:"-"`
}

type ProjectConfig struct {
	OS      string `json:"os"`
	Title   string `json:"title"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Start   string `json:"start"`
	Stop    string `json:"stop"`
	Install string `json:"install"`
}
