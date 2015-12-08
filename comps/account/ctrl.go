package account

import (
	"strings"

	"github.com/go-macaron/captcha"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps/core"
	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

//----------------------------------------------------------
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

//----------------------------------------------------------
// /a/validemail/
func GetValidEmailCtrl(c *macaron.Context) {
	r := core.NewRender(c)

	ve := r.Session.Get("validemail")
	if ve == nil {
		r.RedirectMsg("非法的URL请求,或请求已过期！", "/")
		return
	} else {
		r.Data["email"] = ve.(string)
		r.Session.Delete("validemail")
	}

	u := &models.Users{Email: ve.(string)}
	if !models.NewTr().Read(u, "Email") {
		r.RedirectMsg("不存在此Email注册信息！", "/")
		return
	}

	if u.ValidEmail {
		CleanCookies(c, r.Session)
		r.RedirectMsg("邮箱已通过验证，请返回登录", "/a/signin")
		return
	}

	r.SetCrumb("邮件验证", "/a/validemail/")
	r.AddCss("signup.css")
	r.RHTML(200, "account/valid_email")
}

//----------------------------------------------------------
// /a/signout/
func GetSignoutCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	CleanCookies(c, r.Session)
	r.RedirectMsg("正在准备退出, 请稍候...", "/")
	return
}

//----------------------------------------------------------
// /a/signin/
func GetSigninCtrl(c *macaron.Context) {
	r := core.NewRender(c)

	url := c.QueryEscape("url")
	if r.UserInfo != nil {
		if url != "" {
			c.Resp.Header().Set("Location", url)
		} else {
			c.Redirect("/")
			return
		}
	}

	return_url := url
	if return_url == "" {
		return_url = r.Header().Get("HTTP_REFERER")
	}
	r.Data["return_url"] = return_url

	r.SetCrumb("登录", "/a/signin/")
	r.AddCss("signin.css").AddJs("comps/signin.js")
	r.RHTML(200, "account/signin")
}

//----------------------------------------------------------
// /a/welcomemsg/
func GetWelcomeMsgCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	if _, ok := r.CheckUser(); !ok {
		r.Status(401)
		return
	}

	r.Data["jobs"] = models.AllJobs()
	r.RHTML(200, "account/ajax/welcome_message")
}

//----------------------------------------------------------
// /a/welcometopics/
func GetWelcomeTopicsCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	if _, ok := r.CheckUser(); !ok {
		r.Status(401)
		return
	}

	r.Data["topics"], _ = models.GetTopicExts(8, r.Uid)
	r.RHTML(200, "account/ajax/welcome_get_topics")
}

//----------------------------------------------------------
// /a/welcomeusers/
func GetWelcomeUsersCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	if _, ok := r.CheckUser(); !ok {
		r.Status(401)
		return
	}

	users := ([]*models.Users)(nil)
	if ru := boot.SysSetting.Ra.WelcomeRecmdusers; ru != "" {
		rusers := strings.Split(ru, ",")
		users, _ = models.GetRecommendRandUser(6, r.Uid, rusers)
	} else {
		users, _ = models.GetActivityUsers(6, r.Uid)
	}

	for _, u := range users {
		u.ExtAttrs = map[string]interface{}{"FollowCheck": models.UFollowExistedById(r.Uid, u.Id)}
	}

	r.Data["users"] = users
	r.RHTML(200, "account/ajax/welcome_get_users")
}
