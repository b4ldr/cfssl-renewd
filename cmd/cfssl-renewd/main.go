package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/sign"
	cfssl_config "github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/signer"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/b4ldr/cfssl-renewd/config"
)

var (
	configFile       = flag.String("config", "/etc/cfssl/renewd.yaml", "path to the config file")
	signerConfigFile = flag.String("signer-config", "/etc/cfssl/client-cfssl.conf", "path to the cfssl client config file")
	cfg              config.Config
	cfssl_signer     signer.Signer
)

func main() {
	flag.Parse()
	if err := cleanenv.ReadConfig(*configFile, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var cfssl_cfg cli.Config
	var err error
	cfssl_cfg.CFG, err = cfssl_config.LoadFile(*signerConfigFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	cfssl_signer, err = sign.SignerFromConfig(cfssl_cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for _, request := range cfg.Requests {
		fmt.Printf("%s\n", request)
		request.Set_defaults()
		if err := request.Init_files(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if !(request.Certfile_exist() && request.Key_exist()) {
			if err = request.Gencert(cfssl_signer); err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	}
}
