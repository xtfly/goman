package router

import (
	"github.com/go-macaron/binding"
	"github.com/xtfly/goman/comps/account"
	"github.com/xtfly/goman/comps/home"
	"gopkg.in/macaron.v1"
)

func Route(m *macaron.Macaron) {
	m.Get("/", home.HomeCtrl)

	// account controls
	m.Get("/a/signup/", account.GetSignupCtrl)
	m.Get("/a/signout/", account.GetLogoutCtrl)
	m.Get("/a/validemail/", account.GetValidEmailCtrl)

	// account api
	m.Post("/api/account/check/", account.ApiCheckUserName)
	m.Post("/api/acount/signup/", binding.Bind(account.SignupForm{}), account.ApiUserSignup)
}
