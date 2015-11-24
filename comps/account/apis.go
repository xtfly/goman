package account

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/kits"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

// POST /api/account/check/
func ApiCheckUserName(c *macaron.Context) {
	un := strings.TrimSpace(c.QueryEscape("username"))
	if err, ok := CheckUsernameChar(un); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	if CheckUsernameSensitiveWords(un) || models.UserExistedByName(un) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册"))
		return
	}

	c.JSON(200, comps.NewRestErrResp(1, ""))
}

type SignupForm struct {
	Name      string `form:"user_name"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	Gender    int8   `form:"gender"`
	JobId     string `form:"job_id"`
	Province  string `form:"province"`
	City      string `form:"city"`
	Signature string `form:"signature"`
	CsrfToken string `form:"_csrf"`
	ICode     string `form:"icode"`
	IEmail    string `form:"invitation_email"`
	Agree     bool   `form:"agreement_chk"`
}

// POST /api/account/signup/
func ApiUserSignup(c *macaron.Context, f SignupForm, cpt *captcha.Captcha,
	csrf csrf.CSRF, sf *session.Flash, ss session.Store) {
	if csrf.ValidToken(f.CsrfToken) {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的跨站请求"))
		return
	}

	// if input invalid, will store user_name & email
	defer func() {
		sf.Set("email", f.Email)
		sf.Set("user_name", f.Name)
	}()

	ra := boot.SysSetting.Ra
	switch ra.RegisterType {
	case models.RegTypeClose:
		c.JSON(200, comps.NewRestErrResp(-1, "本站目前关闭注册"))
		return
	case models.RegTypeInvite:
		if f.ICode == "" {
			c.JSON(200, comps.NewRestErrResp(-1, "本站只能通过邀请注册"))
		}
	default:
		break
	}

	if !cpt.VerifyReq(c.Req) {
		c.JSON(200, comps.NewRestErrResp(-1, "请填写正确的验证码"))
		return
	}

	if err, ok := CheckUsernameChar(f.Name); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	var invitation *models.Invitation
	if f.ICode != "" {
		if invitation = models.CheckICodeAvailable(f.ICode); invitation == nil {
			c.JSON(200, comps.NewRestErrResp(-1, "邀请码无效或与邀请邮箱不一致"))
			return
		}
	}

	if CheckUsernameSensitiveWords(f.Name) || models.UserExistedByName(f.Name) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册或包含敏感词或系统保留字"))
		return
	}

	if kits.IsEmail(f.Email) || models.UserExistedByEmail(f.Email) {
		c.JSON(200, comps.NewRestErrResp(-1, "EMail 已经被使用, 或格式不正确"))
		return
	}

	if len(f.Password) < 6 {
		c.JSON(200, comps.NewRestErrResp(-1, "密码长度不符合规则"))
		return
	}

	if !f.Agree {
		c.JSON(200, comps.NewRestErrResp(-1, "你必需同意用户协议才能继续"))
		return
	}

	if f.Gender < models.GenderUnknown || f.Gender > models.GenderFemale {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的性别输入"))
		return
	}

	u := &models.Users{}
	u.UserName = f.Name
	u.Email = f.Email
	u.Password = f.Password
	u.Gender = f.Gender
	u.JobId = com.StrTo(f.JobId).MustInt64()
	u.Province = f.Province
	u.City = f.City
	u.Signature = f.Signature
	u.RegIp = c.RemoteAddr()
	u.Group = &models.UsersGroup{Id: 3} // TODO 未验证会员
	if invitation != nil && f.Email == invitation.Email {
		u.ValidEmail = true
		u.Group = &models.UsersGroup{Id: 4} // TODO 验证会员
	}

	t := models.NewTr()
	uid, ok := u.Add(t)
	if !ok {
		c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
		return
	}
	u.Id = uid

	// 把邀请者加为好友
	if invitation != nil {
		if !models.AddUserFollow(t, uid, invitation.Uid) {
			c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
			return
		}

		if !invitation.Active(t, c.RemoteAddr(), uid) {
			c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
			return
		}
	}

	// 如果不需要email验证
	if !boot.SysSetting.Ra.RegisterValidType || u.Group.Id == 3 || u.ValidEmail {
		SetSigninCookies(c, u)
		c.JSON(200, comps.NewRestRedirectResp("/h/firstlogin_true"))
		return
	}

	ss.Set("validemail", u.Email)
	if !models.NewValidByEmail(t, uid, u.Email) {
		c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
		return
	}

	SetSigninCookies(c, u)
	c.JSON(200, comps.NewRestRedirectResp("/a/validemail/"))
	return

}

//
func SetSigninCookies(c *macaron.Context, u *models.Users) {
	//c.SetCookie("uid", u.Id）
	//c.SetCookie("token", "")
}

//
func CheckUsernameChar(un string) (string, bool) {
	if kits.IsDigit(un) {
		return "用户名不能为纯数字", false
	}

	if strings.ContainsAny(un, "-./") || strings.Contains(un, "__") {
		return "用户名不能包含 - / . % 与连续的下划线", false
	}

	unlen := len(un)
	min := boot.SysSetting.Ra.UsernameLenMin
	max := boot.SysSetting.Ra.UsernameLenMax
	if unlen < min || unlen > max {
		return fmt.Sprintf("用户名长度只能在[%d,%d]", min, max), false
	}

	switch boot.SysSetting.Ra.UsernameRule {
	case models.UserRuleNotLimit:
		break
	case models.UserRuleChineseLetterNumUnline:
		if !kits.IsChineseLetterNumUnline(un) {
			return fmt.Sprintf("请输入大于[%d,%d]字节的用户名, 允许汉字、字母与数字", min, max), false
		}
	case models.UserRuleLetterNumUnline:
		if !kits.IsLetterNumUnline(un) {
			return fmt.Sprintf("请输入[%d,%d]个字母、数字或下划线", min, max), false
		}
	case models.UserRuleChinese:
		if !kits.IsChinese(un) {
			return fmt.Sprintf("请输入[%d,%d]个汉字", min/2, max/2), false
		}
	default:
		break
	}

	return "", true
}

//检查用户名中是否包含敏感词或用户信息保留字
func CheckUsernameSensitiveWords(un string) bool {
	if kits.SensitiveWordExists(un, boot.SysSetting.Cs.SensitiveWords) {
		return true
	}

	cs := boot.SysSetting.Ra.CensorUser
	if len(cs) == 0 {
		return false
	}

	css := strings.Split(cs, "\n")
	for _, c := range css {
		if strings.ToLower(un) == strings.ToLower(strings.TrimSpace(c)) {
			return true
		}
	}

	return false
}
