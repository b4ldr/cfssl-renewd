package request

import (
	"encoding/json"
	"fmt"
	"os"
)

type CSR struct {
	CN    string   `yaml:"cn"`
	Hosts []string `yaml:"hosts"`
	Names []string `yaml:"names"`
	Key   Key      `yaml:"key"`
}

func (c *CSR) json() ([]byte, error) {
	json, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return json, nil
}
func (c *CSR) write(outdir string, perms Perms) error {
	path := fmt.Sprintf("%s/%s.json", outdir, c.CN)
	json, err := c.json()
	if err != nil {
		return err
	}
	os.WriteFile(path, json, perms.mode())
	return nil
}
