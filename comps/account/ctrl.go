package account

import (
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps/core"
	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

// /a/signup
func GetSignupCtrl(c *macaron.Context, cpt *captcha.Captcha, a token.TokenService) {
	r := core.NewRender(c)

	if boot.SysSetting.Ra.SiteClose {
		r.RedirectMsg("本站目前关闭注册", "/")
		return
	}

	icode := c.QueryEscape("icode")
	if boot.SysSetting.Ra.RegisterType == models.RegTypeInvite && icode == "" {
		r.RedirectMsg("本站只接受邀请注册", "/")
		return
	}

	if icode != "" {
		if i := models.CheckICodeAvailable(icode); i != nil {
			r.Data["icode"] = icode
		} else {
			r.RedirectMsg("邀请码无效或已经使用, 请使用新的邀请码", "/")
			return
		}
	}

	r.AddCss("signup.css").AddJs("comps/signup.js")

	r.SetCaptcha(cpt)

	c.Data["jobs"] = models.AllJobs()
	c.Data["csrf_token"], _ = a.GenSysToken(c.RemoteAddr(), 15)

	r.RHTML(200, "account/signup")
}

// /a/validemail/
func GetValidEmailCtrl(c *macaron.Context, ss session.Store) {
	r := core.NewRender(c)

	ve := ss.Get("validemail")
	if ve == nil {
		r.RedirectMsg("非法的URL请求！", "/")
		return
	} else {
		r.Data["email"] = ve.(string)
		ss.Delete("validemail")
	}

	r.AddCss("signup.css")
	r.RHTML(200, "account/valid_email")
}

// /a/signout/
func GetLogoutCtrl(c *macaron.Context, ss session.Store) {
	r := core.NewRender(c)
	CleanCookies(c, ss)
	r.RedirectMsg("正在准备退出, 请稍候...", "/")
	return
}
