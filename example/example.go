package main

import (
	"github.com/cloudfly/confo"
)

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

var (
	logFile    = confo.NewString("log.file", "", "the log file")
	logLevel   = confo.NewBool("log.level", false, "the log level, default is bool")
	serverHost = confo.NewString("server.addr", ":7070", "the address http will serve on")
	serverPort = confo.NewInt("server.port", 8080, "the tcp port will listen on")
)

func main() {
	confo.Parse()

	confo.Usage("log.file")

}
