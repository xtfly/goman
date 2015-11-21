package router

import (
	"github.com/xtfly/goman/comps/account"
	"github.com/xtfly/goman/comps/home"
	"gopkg.in/macaron.v1"
)

func Route(m *macaron.Macaron) {
	m.Get("/", home.HomeCtrl)

	m.Get("/a/signup", account.GetSignupCtrl)

	m.Post("/api/account/check", account.ApiCheckUserName)
}
