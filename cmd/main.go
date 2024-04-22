package main

import (
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../configFiles/", "Path to configuration file")
}

func main() {
}
