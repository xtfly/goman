package core

import (
	"fmt"
	"html"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

type Render struct {
	*macaron.Context
	Uid      int64
	UserInfo *models.Users
}

func NewRender(c *macaron.Context) *Render {
	r := &Render{
		Context: c,
	}
	return r.init()
}

func (c *Render) defDatas() *Render {
	c.Data["X_UA_COMPATIBLE"] = boot.X_UA_COMPATIBLE
	c.Data["VERSION_BUILD"] = boot.VERSION_BUILD
	c.Data["sys"] = boot.SysSetting
	c.Data["uid"] = c.Uid
	return c
}

func (c *Render) defCss() *Render {
	s := boot.SysSetting.Si.Static
	u := boot.SysSetting.Ps.UiStyle
	c.Data["css_files"] = []string{
		s + "/front/css/" + u + "/common.css",
		s + "/front/css/" + u + "/link.css",
	}
	return c
}

func (c *Render) defJs() *Render {
	s := boot.SysSetting.Si.Static
	c.Data["js_files"] = []string{
		s + "/front/js/goman.js",
		s + "/front/js/template.js",
		s + "/front/js/app.js",
	}
	return c
}

func (c *Render) AddCss(css ...string) *Render {
	for _, f := range css {
		c.Data["css_files"] = append(c.Data["css_files"].([]string), boot.SysSetting.Si.Static+"/front/css/"+f)
	}
	return c
}

func (c *Render) AddJs(js ...string) *Render {
	for _, f := range js {
		c.Data["js_files"] = append(c.Data["js_files"].([]string), boot.SysSetting.Si.Static+"/front/js/"+f)
	}
	return c
}

func (c *Render) RHTML(sc int, tn string, data ...interface{}) {
	us := boot.SysSetting.Ps.UiStyle
	c.HTML(sc, us+"/"+tn)
}

func (c *Render) SetCaptcha(cpt *captcha.Captcha) {
	SetCaptcha(c.Context, cpt)
}

func SetCaptcha(c *macaron.Context, cpt *captcha.Captcha) {
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
	c.RHTML(200, "global/show_message")
}

func (c *Render) RedirectMsg(msg, url string) {
	c.RedirectMsgWithDelay(msg, url, 5)
}

type Crumb struct {
	Name string
	Url  string
}

//产生面包屑导航数据,并生成浏览器标题供前端使用
func (c *Render) SetCrumb(name string, url string) {
	sname := html.UnescapeString(name)
	crumbtpl := ([]*Crumb)(nil)
	if crumb := c.Data["crumb"]; crumb == nil {
		crumbtpl = []*Crumb{&Crumb{Name: sname, Url: url}}
	} else {
		crumbtpl = crumb.([]*Crumb)
		crumbtpl = append(crumbtpl, &Crumb{Name: sname, Url: url})
	}
	c.Data["crumb"] = crumbtpl

	var title string
	for _, item := range crumbtpl {
		title = item.Name + "/" + title
	}

	c.Data["page_title"] = html.EscapeString(strings.TrimRight(title, "/"))
}

func (c *Render) init() *Render {
	// 如果登录用户，则加载用户信息
	// c.Data["uid"]是token插件设置
	uid := c.Data["uid"]
	if uid != nil {
		if c.Uid = uid.(int64); c.Uid != 0 {
			c.loadUser()
		}
	}

	c.defDatas().defCss().defJs()

	//
	c.SetCrumb(boot.SysSetting.Si.SiteName, "/")

	return c
}

func (c *Render) loadUser() {
	ss := c.getSessionStore()
	if ss == nil {
		return
	}

	uinfo := ss.Get("uinfo")
	if uinfo == nil {
		c.UserInfo = &models.Users{Id: c.Uid}
		t := models.NewTr()

		if !t.Read(c.UserInfo) {
			return
		}

		// if !t.Read(c.UserInfo.Group) {
		// 	return
		// }

		ss.Set("uinfo", c.UserInfo)
	} else {
		c.UserInfo = uinfo.(*models.Users)
	}
}

func (c *Render) getSessionStore() session.Store {
	ss := (*session.Store)(nil)
	sst := reflect.TypeOf(ss).Elem()
	log.Info(sst)
	ssv := c.GetVal(sst)
	if ssv.CanInterface() {
		return ssv.Interface().(session.Store)
	}
	log.Error("Get session.Store failed.")
	return nil
}
