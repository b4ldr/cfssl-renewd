package request

type Key struct {
	Algo string `yaml:"algo"`
	Size int    `yaml:"size"`
}

func (k *Key) Set_defaults() {
	println("set key defaults")
	if k.Algo == "" {
		k.Algo = "ecdsa"
	}
	if k.Size == 0 {
		k.Size = 256
	}
}
