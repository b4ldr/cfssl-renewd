package request

type Reload struct {
	Services []string `yaml:"services"`
	Command  string   `yaml:"command"`
}
