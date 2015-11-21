package account

import (
	"fmt"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

func GetSignupCtrl(c *macaron.Context, cpt *captcha.Captcha, x csrf.CSRF) {
	comps.DefCxt(c, 0)
	comps.AddCss(c, "signup.css")
	comps.AddJs(c, "comps/signup.js")

	cptvalue, err := cpt.CreateCaptcha()
	if err != nil {
		return
	}
	c.Data["captcha_id"] = cptvalue
	c.Data["captcha_url"] = fmt.Sprintf("%s%s%s.png", cpt.SubURL, cpt.URLPrefix, cptvalue)

	c.Data["jobs"] = models.AllJobs()
	c.Data["csrf_token"] = x.GetToken()

	comps.HTML(c, 200, "account/signup")
}
