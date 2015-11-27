package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/plugins/spider"
	"github.com/xtfly/goman/plugins/token"
	"github.com/xtfly/goman/router"
	"gopkg.in/macaron.v1"
)

func main() {
	log.Debug("Starting server...")

	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(cache.Cacher())
	m.Use(session.Sessioner())
	m.Use(captcha.Captchaer(captcha.Options{Width: 120, Height: 40}))
	m.Use(macaron.Static("static", macaron.StaticOptions{Prefix: "/static"}))
	m.Use(pongo2.Pongoer())
	//m.Use(i18n.I18n(i18n.Options{
	//	Langs: []string{"en-US", "zh-CN"},
	//	Names: []string{"English", "简体中文"},
	//}))
	m.Use(spider.SpiderFunc())
	m.Use(token.Tokener())

	boot.BootStrap()
	router.Route(m)

	m.Run(boot.WebListenIP, boot.WebPort)
}
