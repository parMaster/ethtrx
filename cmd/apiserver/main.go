package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-pkgz/lgr"
	"github.com/parMaster/ethtrx/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		lgr.Fatalf(err.Error())
		os.Exit(1)
	}

	if err := apiserver.Start(config); err != nil {
		lgr.Fatalf("Can't start logserver %s", err.Error())
		os.Exit(1)
	}
}
