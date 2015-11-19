package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/pongo2"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/plugins/auth"
	"github.com/xtfly/goman/plugins/spider"
	"gopkg.in/macaron.v1"
)

func main() {
	log.Debug("Starting server...")

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("static"))
	m.Use(pongo2.Pongoer(pongo2.Options{Directory: "view"}))
	m.Use(i18n.I18n(i18n.Options{
		Langs: []string{"en-US", "zh-CN"},
		Names: []string{"English", "简体中文"},
	}))
	m.Use(spider.SpiderFunc())
	m.Use(auth.Auther())

	boot.BootStrap()
	m.Run(boot.WebListenIP, boot.WebPort)
}
