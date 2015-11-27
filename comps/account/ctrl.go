package account

import (
	"github.com/go-macaron/captcha"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps/core"
	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

func GetSignupCtrl(c *macaron.Context, cpt *captcha.Captcha, a token.TokenService) {
	r := core.NewRender(c)

	if boot.SysSetting.Ra.SiteClose {
		r.RedirectMsg("本站目前关闭注册", "/")
		return
	}

	icode := c.QueryEscape("icode")
	if boot.SysSetting.Ra.RegisterType == models.RegTypeInvite || icode == "" {
		r.RedirectMsg("本站只接受邀请注册", "/")
		return
	}

	if icode != "" {
		if i := models.CheckICodeAvailable(icode); i != nil {
			r.Set("icode", icode)
		} else {
			r.RedirectMsg("邀请码无效或已经使用, 请使用新的邀请码", "/")
			return
		}
	}

	r.AddCss("signup.css").AddJs("comps/signup.js")

	r.SetCaptcha(cpt)

	c.Data["jobs"] = models.AllJobs()
	c.Data["csrf_token"], _ = a.GenSysToken(c.RemoteAddr(), 15)

	r.HTML(200, "account/signup")
}
