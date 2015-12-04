package account

import (
	"strings"
	"time"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

//----------------------------------------------------------
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

//----------------------------------------------------------
// POST /api/account/signup/
func ApiUserSignup(f SignupForm, c *macaron.Context, cpt *captcha.Captcha,
	a token.TokenService, ss session.Store) {
	if !a.ValidToken(c.RemoteAddr(), f.CsrfToken) {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的跨站请求"))
		return
	}

	if !cpt.VerifyReq(c.Req) {
		c.JSON(200, comps.NewRestResp(comps.NewCaptcha(cpt), -1, "请填写正确的验证码"))
		return
	}

	s := NewService()
	u, msg, ok := s.Signup(f, c.RemoteAddr())
	if !ok {
		c.JSON(200, comps.NewRestResp(comps.NewCaptcha(cpt), -1, msg))
		return
	}

	// 如果不需要email验证
	if boot.SysSetting.Ra.RegisterValidType == models.RegValidNone ||
		u.GroupId != models.GroupNotValidated ||
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

//----------------------------------------------------------
// POST /api/account/signin/
func ApiSignin(c *macaron.Context, f SigninForm, a token.TokenService, ss session.Store) {
	u := &models.Users{}
	if !u.CheckSignin(f.Input, f.Password) {
		c.JSON(200, comps.NewRestErrResp(-1, "输入正确的帐号或密码"))
		return
	}

	s := NewService()
	if err, ok := s.CheckSignin(u); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	// 需要审批
	if u.GroupId == models.GroupNotValidated &&
		boot.SysSetting.Ra.RegisterValidType == models.RegValidApproval {
		c.JSON(200, comps.NewRestRedirectResp("/a/validapproval/"))
		return
	}

	//
	u.LastLogin = time.Now()
	u.LastIp = c.RemoteAddr()
	u.LoginCount = u.LoginCount + 1
	if _, ok := models.NewTr().Update(u, "LastLogin", "LastIp", "LoginCount"); !ok {
		// todo log
	}

	CleanCookies(c, ss)
	SetSigninCookies(c, u, a, ss)

	url := ""
	if !u.ValidEmail && boot.SysSetting.Ra.RegisterValidType == models.RegValidEmail {
		ss.Set("validemail", u.Email)
		url = "/a/validemail/"
	} else if u.FirstLogin {
		url = "/h/firstlogin/"
	} else if f.ReturnUrl != "" {
		url = f.ReturnUrl
	}

	c.JSON(200, comps.NewRestRedirectResp(url))
}
