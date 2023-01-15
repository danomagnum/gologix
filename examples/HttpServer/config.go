package main

type PLCConfig struct {
	Name    string
	Address string
	Path    string
}

type ServerConfig struct {
	Address string
	Port    int
}

type AppConfig struct {
	Server ServerConfig
	PLCs   []PLCConfig
}

var Config AppConfig
