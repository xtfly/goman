package core

import (
	"html"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

type Render struct {
	*macaron.Context
	Session  session.Store
	Uid      int64
	UserInfo *UserInfo
}

type UserInfo struct {
	models.Users
	Group *models.UsersGroup
}

func (u *UserInfo) GetAvatar(size string) string {
	avatar := boot.SysSetting.Si.Static + "/common/"
	if u.Avatar == "" {
		switch size {
		case "max":
			avatar = avatar + "avatar-max-img.png"
		case "min":
			avatar = avatar + "avatar-min-img.png"
		case "mid":
			avatar = avatar + "avatar-min-img.png"
		default:
			avatar = avatar + "avatar-img.png"
		}
	} else {
		avatar = u.Avatar
		switch size {
		case "max":
			break
		case "min":
			avatar = strings.Replace(avatar, "max", "min", -1)
		case "mid":
			avatar = strings.Replace(avatar, "max", "mid", -1)
		default:
			break
		}
	}
	return avatar
}

func NewRender(c *macaron.Context) *Render {
	r := &Render{
		Context: c,
	}
	return r.init()
}

func (c *Render) defDatas() *Render {
	c.Data["VERSION_BUILD"] = boot.VERSION_BUILD
	c.Data["sys"] = boot.SysSetting
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
		c.Data["css_files"] = append(c.Data["css_files"].([]string),
			boot.SysSetting.Si.Static+"/front/css/"+boot.SysSetting.Ps.UiStyle+"/"+f)
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
	ci := comps.NewCaptcha(cpt)
	c.Data["captcha_id"] = ci.CaptchaId
	c.Data["captcha_url"] = ci.CaptchaUrl
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

type crumb struct {
	Name string
	Url  string
}

// 产生面包屑导航数据,并生成浏览器标题供前端使用
func (c *Render) SetCrumb(name string, url string) {
	sname := html.UnescapeString(name)
	crumbtpl := ([]*crumb)(nil)
	if scrumb := c.Data["crumb"]; scrumb == nil {
		crumbtpl = []*crumb{&crumb{Name: sname, Url: url}}
	} else {
		crumbtpl = scrumb.([]*crumb)
		crumbtpl = append(crumbtpl, &crumb{Name: sname, Url: url})
	}
	c.Data["crumb"] = crumbtpl

	var title string
	for _, item := range crumbtpl {
		title = item.Name + "/" + title
	}

	c.Data["page_title"] = html.EscapeString(strings.TrimRight(title, "/"))
}

func (c *Render) IsSignin() (string, bool) {
	if c.Uid == 0 || c.UserInfo == nil {
		return "未登录，非法访问", false
	}
	return "", true
}

func (c *Render) init() *Render {
	c.Session = c.sessionStore()

	// 如果登录用户，则加载用户信息
	// c.Data["uid"]是token插件设置
	uid := c.Data["uid"]
	if uid != nil {
		if c.Uid = uid.(int64); c.Uid != 0 {
			c.loadUser()
		}
	}

	c.defDatas().defCss().defJs()
	c.SetCrumb(boot.SysSetting.Si.SiteName, "/")

	return c
}

func (c *Render) loadUser() {
	if c.Session == nil {
		return
	}

	uinfo := c.Session.Get("uinfo")
	if uinfo == nil {
		c.UserInfo = &UserInfo{}
		c.UserInfo.Id = c.Uid
		t := models.NewTr()
		// load the user info from db
		if !t.Read(&c.UserInfo.Users) {
			return
		}

		// load the user group info from db
		c.UserInfo.Group = &models.UsersGroup{Id: c.UserInfo.GroupId}
		if !t.Read(c.UserInfo.Group) {
			return
		}

		c.Session.Set("uinfo", c.UserInfo)
	} else {
		c.UserInfo = uinfo.(*UserInfo)
	}

	c.Data["u"] = c.UserInfo
}

func (c *Render) sessionStore() session.Store {
	ss := (*session.Store)(nil)
	sst := reflect.TypeOf(ss).Elem()
	//log.Info(sst)
	ssv := c.GetVal(sst)
	if ssv.CanInterface() {
		return ssv.Interface().(session.Store)
	}

	log.Error("Get session.Store failed.")
	return nil
}
