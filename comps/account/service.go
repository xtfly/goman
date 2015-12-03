package account

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/kits"
	"github.com/xtfly/goman/models"
)

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

type AccountService struct {
}

func NewService() *AccountService {
	return &AccountService{}
}

//检查用户名中是否包含敏感词或用户信息保留字
func (s *AccountService) CheckUsernameSensitiveWords(un string) bool {
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

// 检查用户名是否合法
func (s *AccountService) CheckUsernameChar(un string) (string, bool) {
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

func (s *AccountService) CheckSignup(f SignupForm) (string, bool) {
	ra := boot.SysSetting.Ra
	switch ra.RegisterType {
	case models.RegTypeClose:
		return "本站目前关闭注册", false
	case models.RegTypeInvite:
		if f.ICode == "" {
			return "本站只能通过邀请注册", false
		}
	default:
		break
	}

	if err, ok := s.CheckUsernameChar(f.Name); !ok {
		return err, ok
	}

	if s.CheckUsernameSensitiveWords(f.Name) || models.UserExistedByName(f.Name) {
		return "用户名已被注册或包含敏感词或系统保留字", false
	}

	if !kits.IsEmail(f.Email) || models.UserExistedByEmail(f.Email) {
		return "EMail 已经被使用, 或格式不正确", false
	}

	if len(f.Password) < 6 {
		return "密码长度不符合规则", false
	}

	if !f.Agree {
		return "你必需同意用户协议才能继续", false
	}

	if f.Gender < models.GenderUnknown || f.Gender > models.GenderFemale {
		return "非法的性别输入", false
	}

	return "", true
}

func (s *AccountService) Signup(f SignupForm, clientip string) (*models.Users, string, bool) {
	if msg, ok := s.CheckSignup(f); !ok {
		return nil, msg, ok
	}

	var invitation *models.Invitation
	if f.ICode != "" {
		if invitation = models.CheckICodeAvailable(f.ICode); invitation == nil {
			return nil, "邀请码无效或与邀请邮箱不一致", false
		}
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
	u.RegIp = clientip
	u.GroupId = models.GroupNotValidated // 未验证会员
	if invitation != nil && f.Email == invitation.Email {
		u.ValidEmail = true
		u.GroupId = models.GroupNormal //  验证会员
	}

	t := models.NewTr().Begin()
	defer t.End()
	uid, ok := u.Add(t)
	if !ok {
		return nil, "内部系统错误", false

	}
	u.Id = uid

	// 把邀请者加为好友
	if invitation != nil {
		if !models.AddUserFollow(t, uid, invitation.Uid) {
			return nil, "内部系统错误", false
		}

		if !invitation.Active(t, clientip, uid) {
			return nil, "内部系统错误", false
		}
	}

	return u, "", true
}

func (s *AccountService) CheckSignin(m *models.Users) (string, bool) {
	if m.Forbidden {
		return "抱歉, 你的账号已经被禁止登录", false
	}

	if boot.SysSetting.Ra.SiteClose &&
		m.GroupId != models.GroupSuperAdmin &&
		m.GroupId != models.GroupWebAdmin {
		return boot.SysSetting.Ra.SiteNotice, false
	}

	return "", true
}

type SigninForm struct {
	Input     string `form:"name"`
	Password  string `form:"password"`
	ReturnUrl string `form:"return_url"`
}
