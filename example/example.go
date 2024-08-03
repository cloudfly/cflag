package main

import "github.com/cloudfly/confo"

type Level string

type Config struct {
	Name string `yaml:"name"`
	Log  struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"log"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		TLS  bool   `yaml:"tls"`
	} `yaml:"server"`
}

func main() {
	var c Config
	err := confo.Load(&c)
	if err != nil {
		panic(err)
	}
}
