package main

import (
	"github.com/fc-testzone/apiserver/net"
	"github.com/fc-testzone/apiserver/utils"
)

type App struct {
	log    *utils.Log
	cfg    *utils.Configs
	server *net.WebServer
}

func NewApp(s *net.WebServer, c *utils.Configs, l *utils.Log) *App {
	return &App{
		server: s,
		cfg:    c,
		log:    l,
	}
}

func (a *App) Start() {
	a.log.SetPath("/var/log/apiserver/")

	var err = a.cfg.LoadFromFile("/etc/apiserver/server.conf")
	if err != nil {
		err = a.cfg.LoadFromFile("server.conf")
		if err != nil {
			a.log.Error("APP", "Fail to load configs", err.Error())
			return
		}
	}
	a.log.Info("APP", "Configs was loaded")

	var webCfg = a.cfg.Settings().Server

	a.log.Info("APP", "Starting web server...")
	err = a.server.Start(webCfg.IP, webCfg.Port)
	if err != nil {
		a.log.Error("APP", "Fail to start web server", err.Error())
	}
}
