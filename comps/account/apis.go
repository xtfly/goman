package account

import (
	"strings"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

// POST /api/account/check/
func ApiCheckUserName(c *macaron.Context) {
	un := strings.TrimSpace(c.QueryEscape("username"))
	s := NewService()
	if err, ok := s.CheckUsernameChar(un); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	if s.CheckUsernameSensitiveWords(un) || models.UserExistedByName(un) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册"))
		return
	}

	c.JSON(200, comps.NewRestErrResp(1, ""))
}

// POST /api/account/signup/
func ApiUserSignup(c *macaron.Context, f SignupForm, cpt *captcha.Captcha,
	a token.TokenService, sf *session.Flash, ss session.Store) {
	if a.ValidToken(c.RemoteAddr(), f.CsrfToken) {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的跨站请求"))
		return
	}

	// if input invalid, will store user_name & email
	defer func() {
		sf.Set("email", f.Email)
		sf.Set("user_name", f.Name)
	}()

	if !cpt.VerifyReq(c.Req) {
		c.JSON(200, comps.NewRestErrResp(-1, "请填写正确的验证码"))
		return
	}

	s := NewService()
	u, msg, ok := s.Signup(f, c.RemoteAddr())
	if !ok {
		c.JSON(200, comps.NewRestErrResp(-1, msg))
		return
	}

	// 如果不需要email验证
	if boot.SysSetting.Ra.RegisterValidType == models.RegValidNone ||
		u.Group.Id == models.GroupNotValidated ||
		u.ValidEmail {
		SetSigninCookies(c, u, a, ss)
		c.JSON(200, comps.NewRestRedirectResp("/h/firstlogin"))
		return
	}

	ss.Set("validemail", u.Email)
	if !models.NewValidByEmail(models.NewTr(), u.Id, u.Email) {
		c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
		return
	}

	SetSigninCookies(c, u, a, ss)
	c.JSON(200, comps.NewRestRedirectResp("/a/validemail/"))
	return
}

// 设置Cookie信息
func SetSigninCookies(c *macaron.Context, u *models.Users, a token.TokenService, ss session.Store) {
	t, _ := a.GenUserToken(c.RemoteAddr(), u.Id, 24*60, token.TokenUser)
	c.SetCookie("utoken", t, 24*60*60)
	ss.Set("utoken", t)
}

func ApiSignin(c *macaron.Context) {

}
