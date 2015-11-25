package comps

import (
	"fmt"

	"github.com/go-macaron/captcha"
	"github.com/xtfly/goman/boot"
	"gopkg.in/macaron.v1"
)

type Render struct {
	*macaron.Context
}

func NewRender(c *macaron.Context) *Render {
	return &Render{
		Context: c,
	}
}

func (c *Render) AddDatas(uid int64) *Render {
	c.Data["X_UA_COMPATIBLE"] = boot.X_UA_COMPATIBLE
	c.Data["VERSION_BUILD"] = boot.VERSION_BUILD
	c.Data["sys"] = boot.SysSetting
	c.Data["uid"] = uid
	return c
}

func (c *Render) HTML(sc int, tn string) {
	us := boot.SysSetting.Ps.UiStyle
	c.HTML(sc, us+"/"+tn)
}

func (c *Render) DefCss() *Render {
	s := boot.SysSetting.Si.Static
	u := boot.SysSetting.Ps.UiStyle
	c.Data["css_files"] = []string{
		s + "/front/css/" + u + "/common.css",
		s + "/front/css/" + u + "/link.css",
	}
	return c
}

func (c *Render) AddCss(css ...string) *Render {
	for _, f := range css {
		c.Data["css_files"] = append(c.Data["css_files"].([]string), boot.SysSetting.Si.Static+"/front/css/"+f)
	}
	return c
}

func (c *Render) DefJs() *Render {
	s := boot.SysSetting.Si.Static
	c.Data["js_files"] = []string{
		s + "/front/js/goman.js",
		s + "/front/js/template.js",
		s + "/front/js/app.js",
	}
	return c
}

func (c *Render) AddJs(js ...string) *Render {
	for _, f := range js {
		c.Data["js_files"] = append(c.Data["js_files"].([]string), boot.SysSetting.Si.Static+"/front/js/"+f)
	}
	return c
}

func (c *Render) DefCxt(uid int64) *Render {
	c.AddDatas(uid).DefCss().DefJs()
	return c
}

func (c *Render) SetCaptcha(cpt *captcha.Captcha) {
	cptvalue, err := cpt.CreateCaptcha()
	if err != nil {
		return
	}

	c.Data["captcha_id"] = cptvalue
	c.Data["captcha_url"] = fmt.Sprintf("%s%s%s.png", cpt.SubURL, cpt.URLPrefix, cptvalue)
}

func (c *Render) RedirectMsgWithDelay(msg, url string, interval int) {
	c.Data["message"] = msg
	c.Data["url_bit"] = url
	c.Data["interval"] = interval
	c.HTML(200, "global/show_message")
}

func (c *Render) RedirectMsg(msg, url string) {
	c.RedirectMsgWithDelay(msg, url, 5)
}
