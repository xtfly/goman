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
	m.Get("/a/signout/", account.GetSignoutCtrl)
	m.Get("/a/signin/", account.GetSigninCtrl)
	m.Get("/a/validemail/", account.GetValidEmailCtrl)
	m.Get("/a/welcomemsg/", account.GetWelcomeMsgCtrl)
	m.Get("/a/welcometopics/", account.GetWelcomeTopicsCtrl)
	m.Get("/a/welcomeusers/", account.GetWelcomeUsersCtrl)

	// account api
	m.Post("/api/account/check/", account.ApiCheckUserName)
	m.Post("/api/account/signup/", binding.Bind(account.SignupForm{}), account.ApiUserSignup)
	m.Post("/api/account/signin/", binding.Bind(account.SigninForm{}), account.ApiSignin)
	m.Post("/api/account/setting/profile", binding.Bind(account.UserSettingForm{}), account.ApiSettingProfile)
	m.Post("/api/account/avatar/upload/", account.ApiUploadAvatar)

	// user home
	m.Get("/h/firstlogin/", home.GetFirstLoginCtrl)
}
