package main

import (
	"os"
	"io/ioutil"
	"syscall"
)

const ENV_CONFIG = "PROMETHEUS_CONFIG"
const ROOT_CONFIG = "/prometheus/etc/default.yml"
const ROOT_BIN = "/usr/local/bin"
const BIN_PROMETHEUS = "prometheus"

func main() {
	envToFile()
	exec()
}

func envToFile(){
	if config, ok := os.LookupEnv(ENV_CONFIG); ok {
		err := ioutil.WriteFile(ROOT_CONFIG, []byte(config), os.ModePerm)
		if err != nil {
			os.Exit(1)
		}
	}
}

func exec(){
	if err := syscall.Exec(ROOT_BIN + "/" + BIN_PROMETHEUS, []string{BIN_PROMETHEUS, "--config.file", ROOT_CONFIG, "--web.listen-address=0.0.0.0:3000", "--log.format=json", "--auto-gomaxprocs", "--auto-gomemlimit", "--storage.tsdb.path=/prometheus/var"}, os.Environ()); err != nil {
		os.Exit(1)
	}
}