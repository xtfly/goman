package account

import (
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

func GetSignupCtrl(c *macaron.Context, cpt *captcha.Captcha, x csrf.CSRF) {
	r := comps.NewRender(c)
	r.DefCxt(0).AddCss("signup.css").AddJs("comps/signup.js")

	r.SetCaptcha(cpt)

	c.Data["jobs"] = models.AllJobs()
	c.Data["csrf_token"] = x.GetToken()

	r.HTML(200, "account/signup")
}
